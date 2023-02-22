@api @skipOnOcV10
Feature: Share spaces
  As the owner of a space
  I want to be able to add members to a space, and to remove access for them

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "share space" with the default quota using the GraphApi
    And using spaces DAV path


  Scenario Outline: A Space Admin can share a space to another user
    When user "Alice" shares a space "share space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And the user "Brian" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | id        | %space_id%  |
      | name      | share space |
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario: A user can see who has been granted access
    When user "Alice" shares a space "share space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    Then the user "Alice" should have a space called "share space" granted to user "Brian" with role "viewer"


  Scenario: A user can see that the group has been granted access
    Given group "sales" has been created
    When user "Alice" shares a space "share space" with settings:
      | shareWith | sales  |
      | shareType | 8      |
      | role      | viewer |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the user "Alice" should have a space called "share space" granted to group "sales" with role "viewer"


  Scenario: A user can see a file in a received shared space
    Given user "Alice" has uploaded a file inside space "share space" with content "Test" to "test.txt"
    And user "Alice" has created a folder "Folder Main" in space "share space"
    When user "Alice" shares a space "share space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    Then for user "Brian" the space "share space" should contain these entries:
      | test.txt    |
      | Folder Main |


  Scenario: When a user unshares a space, the space becomes unavailable to the receiver
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    And the user "Brian" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | id        | %space_id%  |
      | name      | share space |
    When user "Alice" unshares a space "share space" to user "Brian"
    Then the HTTP status code should be "200"
    And the user "Brian" should not have a space called "share space"


  Scenario Outline: Owner of a space cannot see the space after removing his access to the space
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith | Brian   |
      | role      | manager |
    When user "<user>" unshares a space "share space" to user "Alice"
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" owned by "Alice" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | id        | %space_id%  |
      | name      | share space |
    But the user "Alice" should not have a space called "share space"
    Examples:
      | user  |
      | Alice |
      | Brian |


  Scenario: A user can add another user to the space managers to enable him
    Given user "Alice" has uploaded a file inside space "share space" with content "Test" to "test.txt"
    When user "Alice" shares a space "share space" with settings:
      | shareWith | Brian   |
      | role      | manager |
    Then the user "Brian" should have a space called "share space" granted to "Brian" with role "manager"
    When user "Brian" shares a space "share space" with settings:
      | shareWith | Bob    |
      | role      | viewer |
    Then the user "Bob" should have a space called "share space" granted to "Bob" with role "viewer"
    And for user "Bob" the space "share space" should contain these entries:
      | test.txt |


  Scenario Outline: A user cannot share a disabled space to another user
    Given user "Alice" has disabled a space "share space"
    When user "Alice" shares a space "share space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    Then the HTTP status code should be "404"
    And the OCS status code should be "404"
    And the OCS status message should be "Wrong path, file/folder doesn't exist"
    And the user "Brian" should not have a space called "share space"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: A user with manager role can share a space to another user
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith | Brian   |
      | role      | manager |
    When user "Brian" shares a space "share space" with settings:
      | shareWith | Bob    |
      | role      | <role> |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And the user "Bob" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | id        | %space_id%  |
      | name      | share space |
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: A user with editor or viewer role cannot share a space to another user
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" shares a space "share space" with settings:
      | shareWith | Bob        |
      | role      | <new_role> |
    Then the HTTP status code should be "404"
    And the OCS status code should be "404"
    And the OCS status message should be "No share permission"
    And the user "Bob" should not have a space called "share space"
    Examples:
      | role   | new_role |
      | editor | manager  |
      | editor | editor   |
      | editor | viewer   |
      | viewer | manager  |
      | viewer | editor   |
      | viewer | viewer   |


  Scenario Outline: space manager can change the role of space members
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Alice" shares a space "share space" with settings:
      | shareWith | Brian      |
      | role      | <new_role> |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the user "Alice" should have a space called "share space" granted to "Brian" with role "<new_role>"
    Examples:
      | role    | new_role |
      | editor  | manager  |
      | editor  | viewer   |
      | viewer  | manager  |
      | viewer  | editor   |
      | manager | editor   |
      | manager | viewer   |


  Scenario Outline: user without manager role cannot change the role of space members
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    And user "Alice" has shared a space "share space" with settings:
      | shareWith | Bob    |
      | role      | viewer |
    When user "Brian" updates the space "share space" with settings:
      | shareWith | Bob        |
      | role      | <new_role> |
    Then the HTTP status code should be "404"
    And the OCS status code should be "404"
    And the user "Alice" should have a space called "share space" granted to "Bob" with role "viewer"
    Examples:
      | role   | new_role |
      | editor | manager  |
      | editor | viewer   |
      | viewer | manager  |
      | viewer | editor   |


  Scenario Outline: A user shares a space with a group
    Given group "group2" has been created
    And the administrator has added a user "Brian" to the group "group2" using GraphApi
    And the administrator has added a user "Bob" to the group "group2" using GraphApi
    When user "Alice" shares a space "share space" with settings:
      | shareWith | group2 |
      | shareType | 8      |
      | role      | <role> |
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | name      | share space |
    And the user "Bob" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | name      | share space |
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: The user has no access to the space if access for the group has been removed
    Given group "group2" has been created
    And the administrator has added a user "Brian" to the group "group2" using GraphApi
    And user "Alice" has shared a space "share space" with settings:
      | shareWith | group2 |
      | shareType | 8      |
      | role      | <role> |
    When user "Alice" unshares a space "share space" to group "group2"
    Then the HTTP status code should be "200"
    And the user "Brian" should not have a space called "share space"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario: The user has no access to the space if he has been removed from the group
    Given group "group2" has been created
    And the administrator has added a user "Brian" to the group "group2" using GraphApi
    And the administrator has added a user "Bob" to the group "group2" using GraphApi
    And user "Alice" has shared a space "share space" with settings:
      | shareWith | group2 |
      | shareType | 8      |
      | role      | editor |
    When the administrator removes the following users from the following groups using the Graph API
      | username | groupname |
      | Brian    | group2    |
    Then the HTTP status code of responses on all endpoints should be "204"
    And the user "Brian" should not have a space called "share space"
    But the user "Bob" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | name      | share space |


  Scenario: Users don't have access to the space if the group has been deleted
    Given group "group2" has been created
    And the administrator has added a user "Brian" to the group "group2" using GraphApi
    And the administrator has added a user "Bob" to the group "group2" using GraphApi
    And user "Alice" has shared a space "share space" with settings:
      | shareWith | group2 |
      | shareType | 8      |
      | role      | editor |
    When the administrator deletes group "group2" using the Graph API
    Then the HTTP status code should be "204"
    And the user "Brian" should not have a space called "share space"
    And the user "Bob" should not have a space called "share space"


  Scenario: User increases permissions for one member of the group or for the entire group
    Given group "sales" has been created
    And the administrator has added a user "Brian" to the group "sales" using GraphApi
    And user "Alice" has shared a space "share space" with settings:
      | shareWith | sales  |
      | shareType | 8      |
      | role      | viewer |
    When user "Brian" uploads a file inside space "share space" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    When user "Alice" shares a space "share space" with settings:
      | shareWith | Brian  |
      | role      | editor |
    Then the HTTP status code should be "200"
    When user "Brian" uploads a file inside space "share space" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"


  Scenario: User increases permissions for the group, so the user's permissions are increased
    Given group "sales" has been created
    And the administrator has added a user "Brian" to the group "sales" using GraphApi
    And user "Alice" has shared a space "share space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    When user "Brian" uploads a file inside space "share space" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    When user "Alice" shares a space "share space" with settings:
      | shareWith | sales  |
      | shareType | 8      |
      | role      | editor |
    Then the HTTP status code should be "200"
    When user "Brian" uploads a file inside space "share space" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"


  Scenario Outline: A Space Admin can share a space to the user with an expiration date
    When user "Alice" shares a space "share space" with settings:
      | shareWith  | Brian                    |
      | role       | <role>                   |
      | expireDate | 2042-03-25T23:59:59+0100 |
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" granted to user "Brian" with role "<role>" and expiration date "2042-03-25"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: A Space Admin can share a space to the group with an expiration date
    Given group "sales" has been created
    And the administrator has added a user "Brian" to the group "sales" using GraphApi
    When user "Alice" shares a space "share space" with settings:
      | shareWith  | sales                    |
      | shareType  | 8                        |
      | role       | <role>                   |
      | expireDate | 2042-03-25T23:59:59+0100 |
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" granted to group "sales" with role "<role>" and expiration date "2042-03-25"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: update the expiration date of a space in user share
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith  | Brian                    |
      | role       | <role>                   |
      | expireDate | 2042-03-25T23:59:59+0100 |
    When user "Alice" updates the space "share space" with settings:
      | shareWith  | Brian                         |
      | expireDate | 2044-01-01T23:59:59.999+01:00 |
      | role       | <role>                        |
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" granted to user "Brian" with role "<role>" and expiration date "2044-01-01"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: update the expiration date of a space in group share
    Given group "sales" has been created
    And the administrator has added a user "Brian" to the group "sales" using GraphApi
    And user "Alice" has shared a space "share space" with settings:
      | shareWith  | sales                    |
      | shareType  | 8                        |
      | role       | <role>                   |
      | expireDate | 2042-03-25T23:59:59+0100 |
    When user "Alice" updates the space "share space" with settings:
      | shareWith  | sales                         |
      | shareType  | 8                             |
      | expireDate | 2044-01-01T23:59:59.999+01:00 |
      | role       | <role>                        |
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" granted to group "sales" with role "<role>" and expiration date "2044-01-01"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: delete the expiration date of a space in user share
    Given user "Alice" has shared a space "share space" with settings:
      | shareWith  | Brian                    |
      | role       | <role>                   |
      | expireDate | 2042-03-25T23:59:59+0100 |
    When user "Alice" updates the space "share space" with settings:
      | shareWith  | Brian  |
      | expireDate |        |
      | role       | <role> |
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" granted to user "Brian" with role "<role>" and expiration date ""
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: update the expiration date of a space in group share
    Given group "sales" has been created
    And the administrator has added a user "Brian" to the group "sales" using GraphApi
    And user "Alice" has shared a space "share space" with settings:
      | shareWith  | sales                    |
      | shareType  | 8                        |
      | role       | <role>                   |
      | expireDate | 2042-03-25T23:59:59+0100 |
    When user "Alice" updates the space "share space" with settings:
      | shareWith  | sales  |
      | shareType  | 8      |
      | expireDate |        |
      | role       | <role> |
    Then the HTTP status code should be "200"
    And the user "Brian" should have a space called "share space" granted to group "sales" with role "<role>" and expiration date ""
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |    
