@api @skipOnOcV10
Feature: get groups and their members
  As an admin
  I want to be able to get groups
  So that I can see all the groups and their members

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario: admin user lists all the groups
    Given group "tea-lover" has been created
    And group "coffee-lover" has been created
    And group "h2o-lover" has been created
    When user "Alice" gets all the groups using the Graph API
    Then the HTTP status code should be "200"
    And the extra groups returned by the API should be
      | tea-lover    |
      | coffee-lover |
      | h2o-lover    |


  Scenario Outline: normal user cannot get the groups list
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "<role>" using the settings api
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And group "h2o-lover" has been created
    When user "Brian" gets all the groups using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |


  Scenario: admin user gets users of a group
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
    And group "tea-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    And user "Carol" has been added to group "tea-lover"
    When user "Alice" gets all the members of group "tea-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the users returned by the API should be
      | Brian |
      | Carol |


  Scenario Outline: normal user tries to get users of a group
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "<role>" using the settings api
    And group "tea-lover" has been created
    When user "Brian" gets all the members of group "tea-lover" using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |


  Scenario: admin user gets all groups along with its member's information
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    And user "Carol" has been added to group "tea-lover"
    When user "Alice" retrieves all groups along with their members using the Graph API
    Then the HTTP status code should be "200"
    And the group 'coffee-lover' should have the following member information
      | displayName  | id        | mail              | onPremisesSamAccountName |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    |
    And the group 'tea-lover' should have the following member information
      | displayName  | id        | mail              | onPremisesSamAccountName |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    |
      | Carol King   | %uuid_v4% | carol@example.org | Carol                    |


  Scenario Outline: normal user gets all groups along with their members information
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "<role>" using the settings api
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    When user "Brian" retrieves all groups along with their members using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |


  Scenario: admin user gets a group along with their members information
    Given user "Brian" has been created with default attributes and without skeleton files
    And group "tea-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "tea-lover"
    When user "Alice" gets all the members information of group "tea-lover" using the Graph API
    And the group 'tea-lover' should have the following member information
      | displayName  | id        | mail              | onPremisesSamAccountName |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    |

  @issue-5604
  Scenario Outline: normal user gets a group along with their members information
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "<role>" using the settings api
    And group "tea-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "tea-lover"
    When user "Brian" gets all the members information of group "tea-lover" using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |


  Scenario: Get details of a group
    Given group "tea-lover" has been created
    When user "Alice" gets details of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "displayName",
        "id"
      ],
      "properties": {
        "displayName": {
          "type": "string",
          "enum": ["tea-lover"]
        },
        "id": {
          "type": "string",
          "pattern": "^%group_id_pattern%$"
        }
      }
    }
    """


  Scenario Outline: Get details of group with UTF-8 characters name
    Given group "<group>" has been created
    When user "Alice" gets details of the group "<group>" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "displayName",
        "id"
      ],
      "properties": {
        "displayName": {
          "type": "string",
          "enum": ["<group>"]
        },
        "id": {
          "type": "string",
          "pattern": "^%group_id_pattern%$"
        }
      }
    }
    """
    Examples:
      | group           |
      | España§àôœ€     |
      | नेपाली            |
      | $x<=>[y*z^2+1]! |
      | եòɴԪ˯ΗՐΛɔπ     |


  Scenario: admin user tries to get group information of non-existing group
    When user "Alice" gets details of the group "non-existing" using the Graph API
    Then the HTTP status code should be "404"
