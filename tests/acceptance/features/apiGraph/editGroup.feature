Feature: edit group name
  As an admin
  I want to be able to edit group name
  So that I can manage group name

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API

  @issue-5977
  Scenario Outline: admin user renames a group
    Given group "<old_group>" has been created
    When user "Alice" renames group "<old_group>" to "<new_group>" using the Graph API
    Then the HTTP status code should be "204"
    And group "<old_group>" should not exist
    And group "<new_group>" should exist
    Examples:
      | old_group | new_group     |
      | grp1      | grp101        |
      | grp1      | España§àôœ€   |
      | grp1      | नेपाली          |
      | grp1      | $x<=>[y*z^2]! |
      | grp1      | staff?group   |
      | grp1      | 50%pass       |

  @issue-5938
  Scenario Outline: user other than the admin can't rename a group
    Given the administrator has assigned the role "<role>" to user "Alice" using the Graph API
    And group "grp1" has been created
    When user "Alice" tries to rename group "grp1" to "grp101" using the Graph API
    Then the HTTP status code should be "403"
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: admin user tries to rename nonexistent group
    When user "Alice" tries to rename a nonexistent group to "grp1" using the Graph API
    Then the HTTP status code should be "404"
    And group "grp1" should not exist
