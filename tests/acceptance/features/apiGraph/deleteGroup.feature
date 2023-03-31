@api @skipOnOcV10
Feature: delete groups
  As an admin
  I want to be able to delete groups
  So that I can remove unnecessary groups

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario Outline: admin user deletes a group
    Given group "<group_id>" has been created
    When user "Alice" deletes group "<group_id>" using the Graph API
    Then the HTTP status code should be "204"
    And group "<group_id>" should not exist
    Examples:
      | group_id            | comment                               |
      | simplegroup         | nothing special here                  |
      | España§àôœ€         | special European and other characters |
      | नेपाली                | Unicode group name                    |
      | brand-new-group     | dash                                  |
      | the.group           | dot                                   |
      | left,right          | comma                                 |
      | 0                   | The "false" group                     |
      | Finance (NP)        | Space and brackets                    |
      | Admin&Finance       | Ampersand                             |
      | admin:Pokhara@Nepal | Colon and @                           |
      | maint+eng           | Plus sign                             |
      | $x<=>[y*z^2]!       | Maths symbols                         |
      | Mgmt\Middle         | Backslash                             |
      | 😁 😂               | emoji                                 |
      | maintenance#123     | Hash sign                             |
      | 50%25=0             | %25 literal looks like an escaped "%" |
      | staff?group         | Question mark                         |
      | Mgmt/Sydney         | Slash (special escaping happens)      |
      | Mgmt//NSW/Sydney    | Multiple slash                        |
      | priv/subadmins/1    | Subadmins mentioned not at the end    |
      | var/../etc          | using slash-dot-dot                   |

  @issue-5083
  Scenario Outline: admin user deletes a group having % (as only special char) in its name
    Given group "<group_id>" has been created
    When user "Alice" deletes group "<group_id>" using the Graph API
    Then the HTTP status code should be "204"
    And group "<group_id>" should not exist
    Examples:
      | group_id            | comment                                 |
      | 50%pass             | Percent sign (special escaping happens) |
      | 50%2Eagle           | %2E literal looks like an escaped "."   |
      | 50%2Fix             | %2F literal looks like an escaped slash |

  @issue-5938
  Scenario Outline: user other than the admin can't delete a group
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "<role>" using the settings api
    And group "new-group" has been created
    When user "Brian" tries to delete group "new-group" using the Graph API
    Then the HTTP status code should be "403"
    And group "new-group" should exist
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |

  @issue-903
  Scenario: deleted group should not be listed in the sharees list
    Given group "grp1" has been created
    And group "grp2" has been created
    And user "Alice" has uploaded file with content "sample text" to "lorem.txt"
    And user "Alice" has shared file "lorem.txt" with group "grp1"
    And user "Alice" has shared file "lorem.txt" with group "grp2"
    And group "grp1" has been deleted
    When user "Alice" gets all the shares from the file "lorem.txt" using the sharing API
    Then the HTTP status code should be "200"
    And group "grp2" should be included in the response
    But group "grp1" should not be included in the response


  Scenario: user should not see share received via deleted group
    Given user "Alice" has uploaded file with content "sample text" to "lorem.txt"
    And user "Brian" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has shared file "lorem.txt" with group "grp1"
    And user "Brian" has accepted share "/lorem.txt" offered by user "Alice"
    When user "Alice" deletes group "grp1" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should not have any received shares
