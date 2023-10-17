Feature: change shared resource
  As a user
  I want to change the shared resource
  So that I can organize them as per my necessity

  Background:
    Given using spaces DAV path
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @issue-4421
  Scenario: move files between shares by different users
    Given user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has created folder "/PARENT"
    And user "Brian" has created folder "/PARENT"
    And user "Alice" has moved file "textfile0.txt" to "PARENT/from_alice.txt" in space "Personal"
    And user "Alice" has shared folder "/PARENT" with user "Carol"
    And user "Brian" has shared folder "/PARENT" with user "Carol"
    When user "Carol" moves file "PARENT/from_alice.txt" to "PARENT (1)/from_alice.txt" in space "Shares" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Carol" folder "PARENT" of the space "Shares" should not contain these entries:
      | from_alice.txt |
    And for user "Carol" folder "PARENT (1)" of the space "Shares" should contain these entries:
      | from_alice.txt |


  Scenario Outline: overwrite a received file share
    Given the administrator has assigned the role "<userRole>" to user "Brian" using the Graph API
    And user "Alice" has uploaded file with content "old content version 1" to "/textfile1.txt"
    And user "Alice" has uploaded file with content "old content version 2" to "/textfile1.txt"
    And user "Alice" has shared file "/textfile1.txt" with user "Brian"
    When user "Brian" uploads a file inside space "Shares" with content "this is a new content" to "textfile1.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should contain these entries:
      | textfile1.txt |
    And for user "Brian" the content of the file "/textfile1.txt" of the space "Shares" should be "this is a new content"
    And for user "Alice" the content of the file "/textfile1.txt" of the space "Personal" should be "this is a new content"
    When user "Alice" downloads version of the file "/textfile1.txt" with the index "2" of the space "Personal" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded content should be "old content version 1"
    When user "Alice" downloads version of the file "/textfile1.txt" with the index "1" of the space "Personal" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded content should be "old content version 2"
    When user "Brian" tries to get version of the file "/textfile1.txt" with the index "1" of the space "Shares" using the WebDAV API
    Then the HTTP status code should be "403"
    Examples:
      | userRole    |
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |
