@api @skipOnOcV10
Feature: Restoring space
  As a manager of space
  I want to be able to restore a disabled space.
  Only manager can restore disabled space
  The restored space must be visible to the other participants without loss of data

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "restore a space" of type "project" with quota "10"
    And using spaces DAV path


  Scenario: An owner can restore a Space via the Graph API
    Given user "Alice" has disabled a space "restore a space"
    When user "Alice" restores a disabled space "restore a space"
    Then the HTTP status code should be "200"


  Scenario: Participants can see the data after the space is restored
    Given user "Alice" has created a folder "mainFolder" in space "restore a space"
    And user "Alice" has uploaded a file inside space "restore a space" with content "example" to "test.txt"
    And user "Alice" has shared a space "restore a space" with settings:
      | shareWith | Brian  |
      | role      | editor |
    And user "Alice" has shared a space "restore a space" with settings:
      | shareWith | Bob    |
      | role      | viewer |
    And user "Alice" has disabled a space "restore a space"
    When user "Alice" restores a disabled space "restore a space"
    Then for user "Alice" the space "restore a space" should contain these entries:
      | test.txt   |
      | mainFolder |
    And for user "Brian" the space "restore a space" should contain these entries:
      | test.txt   |
      | mainFolder |
    And for user "Bob" the space "restore a space" should contain these entries:
      | test.txt   |
      | mainFolder |


  Scenario: Participant can create data in the space after restoring
    Given user "Alice" has shared a space "restore a space" with settings:
      | shareWith | Brian  |
      | role      | editor |
    And user "Alice" has disabled a space "restore a space"
    And user "Alice" has restored a disabled space "restore a space"
    When user "Brian" creates a folder "mainFolder" in space "restore a space" using the WebDav Api
    And user "Brian" uploads a file inside space "restore a space" with content "test" to "test.txt" using the WebDAV API
    Then for user "Brian" the space "restore a space" should contain these entries:
      | test.txt   |
      | mainFolder |


  Scenario Outline: User without space manager role cannot restore space
    Given user "Alice" has shared a space "restore a space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    And user "Alice" has disabled a space "restore a space"
    When user "Brian" tries to restore a disabled space "restore a space" owned by user "Alice"
    Then the HTTP status code should be "404"
    Examples:
      | role   |
      | viewer |
      | editor |


  Scenario Outline: User with role User and Guest cannot restore space
    Given the administrator has given "Brian" the role "<role>" using the settings api
    And user "Alice" has disabled a space "restore a space"
    When user "Brian" tries to restore a disabled space "restore a space" owned by user "Alice"
    Then the HTTP status code should be "404"
    Examples:
      | role  |
      | User  |
      | Guest |


  Scenario Outline: Admin and Space Admin can restore others space
    Given the administrator has given "Brian" the role "<role>" using the settings api
    And user "Alice" has disabled a space "restore a space"
    When user "Brian" tries to restore a disabled space "restore a space" owned by user "Alice"
    Then the HTTP status code should be "200"
    Examples:
      | role        |
      | Admin       |
      | Space Admin |
