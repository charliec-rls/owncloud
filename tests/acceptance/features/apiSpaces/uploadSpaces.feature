@api
Feature: Upload files into a space
  As a user
  I want to be able to create folders and files in the space
  So that I can store various information in them

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project Ceres" of type "project" with quota "2000"
    And using spaces DAV path


  Scenario Outline: user creates a folder in the space via the Graph API
    Given user "Alice" has shared a space "Project Ceres" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" creates a folder "mainFolder" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "<code>"
    And for user "Brian" the space "Project Ceres" <shouldOrNot> contain these entries:
      | mainFolder |
    Examples:
      | role    | code | shouldOrNot |
      | manager | 201  | should      |
      | editor  | 201  | should      |
      | viewer  | 403  | should not  |


  Scenario Outline: user uploads a file in shared space via the Graph API
    Given user "Alice" has shared a space "Project Ceres" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" uploads a file inside space "Project Ceres" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "<code>"
    And for user "Brian" the space "Project Ceres" <shouldOrNot> contain these entries:
      | test.txt |
    Examples:
      | role    | code | shouldOrNot |
      | manager | 201  | should      |
      | editor  | 201  | should      |
      | viewer  | 403  | should not  |


  Scenario: user can create subfolders in a space via the Graph API
    When user "Alice" creates a subfolder "mainFolder/subFolder1/subFolder2" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Ceres" should contain these entries:
      | mainFolder |
    And for user "Alice" folder "mainFolder/subFolder1/" of the space "Project Ceres" should contain these entries:
      | subFolder2 |


  Scenario: user can create a folder and upload a file to a space
    When user "Alice" creates a folder "NewFolder" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "201"
    And user "Alice" uploads a file inside space "Project Ceres" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Ceres" should contain these entries:
      | NewFolder |
      | test.txt  |


  Scenario: user cannot create a folder or a file in a space if they do not have permission
    When user "Bob" creates a folder "forAlice" in space "Project Ceres" owned by the user "Alice" using the WebDav Api
    Then the HTTP status code should be "404"
    When user "Bob" uploads a file inside space "Project Ceres" owned by the user "Alice" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "404"
    And for user "Alice" the space "Project Ceres" should not contain these entries:
      | forAlice |
      | test.txt |


  Scenario: user cannot create folder with an existing name
    Given user "Alice" has created a folder "NewFolder" in space "Project Ceres"
    When user "Alice" creates a folder "NewFolder" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "405"


  Scenario Outline: user cannot create subfolder in a nonexistent folder
    When user "Alice" tries to create subfolder "<path>" in a nonexistent folder of the space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "409"
    Examples:
      | path        |
      | foo/bar     |
      | foo/bar/baz |
