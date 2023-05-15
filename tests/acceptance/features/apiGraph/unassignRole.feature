@api
Feature: unassign user role
  As an admin
  I want to unassign the role of user
  So that the role of user is set to default

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: admin user unassigns the role of a user with different role using graph API
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the administrator has assigned the role "<role>" to user "Brian" using the Graph API
    When user "Alice" unassigns the role of user "Brian" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should not have any role assigned
    When user "Brian" tries to create a group "assigns-role-to-default" using the Graph API
    And user "Brian" should have the role "User" assigned
    Examples:
      | role        |
      | Admin       |
      | Space Admin |
      | User        |
      | Guest       |

  @issue-6035
  Scenario: admin user tries to unassign his/her own role using graph API
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When user "Alice" tries to unassign the role of user "Alice" using the Graph API
    Then the HTTP status code should be "500"
    And user "Alice" should have the role "Admin" assigned
