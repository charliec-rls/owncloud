@api @skipOnOcV10
Feature: get users
  As an admin
  I want to be able to retrieve user information
  So that I can see the information

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @skipOnStable2.0
  Scenario: admin user gets the information of a user
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |

  @issue-5125
  Scenario Outline: non-admin user tries to get the information of a user
    Given the administrator has given "Alice" the role "<role>" using the settings api
    And the administrator has given "Brian" the role "<userRole>" using the settings api
    When user "Brian" tries to get information of user "Alice" using Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | userRole    | role        |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | Guest       |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | Guest       |
      | User        | Admin       |
      | Guest       | Space Admin |
      | Guest       | User        |
      | Guest       | Guest       |
      | Guest       | Admin       |

  @skipOnStable2.0
  Scenario: admin user gets all users
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When user "Alice" gets all users using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain following users with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    | true           |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |

  @skipOnStable2.0
  Scenario: admin user gets all users include disabled users
    Given the administrator has given "Alice" the role "Admin" using the settings api
    And the user "Alice" has disabled user "Brian" using the Graph API
    When user "Alice" gets all users using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain following users with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    | true           |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | false          |


  Scenario Outline: non-admin user tries to get all users
    Given the administrator has given "Alice" the role "<userRole>" using the settings api
    When user "Brian" tries to get all users using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | userRole    |
      | Space Admin |
      | User        |
      | Guest       |

  @skipOnStable2.0
  Scenario: admin user gets the drive information of a user
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When the user "Alice" gets user "Brian" along with his drive information using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |
    And the user retrieve API response should contain the following drive information:
      | driveType         | personal                         |
      | driveAlias        | personal/brian                   |
      | id                | %space_id%                       |
      | name              | Brian Murphy                     |
      | owner@@@user@@@id | %user_id%                        |
      | quota@@@state     | normal                           |
      | root@@@id         | %space_id%                       |
      | root@@@webDavUrl  | %base_url%/dav/spaces/%space_id% |
      | webUrl            | %base_url%/f/%space_id%          |

  @skipOnStable2.0
  Scenario Outline: non-admin user gets his/her own drive information
    Given the administrator has given "Brian" the role "<userRole>" using the settings api
    When the user "Brian" gets his drive information using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |
    And the user retrieve API response should contain the following drive information:
      | driveType         | personal                         |
      | driveAlias        | personal/brian                   |
      | id                | %space_id%                       |
      | name              | Brian Murphy                     |
      | owner@@@user@@@id | %user_id%                        |
      | quota@@@state     | normal                           |
      | root@@@id         | %space_id%                       |
      | root@@@webDavUrl  | %base_url%/dav/spaces/%space_id% |
      | webUrl            | %base_url%/f/%space_id%          |
    Examples:
      | userRole    |
      | Space Admin |
      | User        |
      | Guest       |

  @skipOnStable2.0
  Scenario: admin user gets the group information of a user
    Given the administrator has given "Alice" the role "Admin" using the settings api
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets user "Brian" along with his group information using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled | memberOf                |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           | tea-lover, coffee-lover |

  @issue-5125
  Scenario Outline: non-admin user tries to get the group information of a user
    Given the administrator has given "Alice" the role "<userRole>" using the settings api
    And the administrator has given "Brian" the role "<role>" using the settings api
    And group "coffee-lover" has been created
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets user "Brian" along with his group information using Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | userRole    | role        |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | Guest       |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | Guest       |
      | User        | Admin       |
      | Guest       | Space Admin |
      | Guest       | User        |
      | Guest       | Guest       |
      | Guest       | Admin       |

  @skipOnStable2.0
  Scenario: admin user gets all users of certain groups
    Given the administrator has given "Alice" the role "Admin" using the settings api
    And user "Carol" has been created with default attributes and without skeleton files
    And the user "Alice" has disabled user "Carol" using the Graph API
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets all users of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain following users with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    | true           |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |
    But the API response should not contain following user with the information:
      | displayName | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Carol King  | %uuid_v4% | carol@example.org | Carol                    | false          |
    When the user "Alice" gets all users of two groups "tea-lover,coffee-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain following user with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |
    But the API response should not contain following users with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    | true           |
      | Carol King   | %uuid_v4% | carol@example.org | Carol                    | true           |

  @skipOnStable2.0
  Scenario: admin user gets all users of certain groups
    Given the administrator has given "Alice" the role "Admin" using the settings api
    And user "Carol" has been created with default attributes and without skeleton files
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And group "wine-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    And user "Carol" has been added to group "wine-lover"
    When the user "Alice" gets all users from that are members in the group "tea-lover" or the group "coffee-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain following users with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    | true           |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |
    But the API response should not contain following user with the information:
      | displayName | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Carol King  | %uuid_v4% | carol@example.org | Carol                    | false          |

  @skipOnStable2.0
  Scenario Outline: non admin user tries to get users of certain groups
    Given the administrator has given "Alice" the role "Admin" using the settings api
    And the administrator has given "Brian" the role "<role>" using the settings api
    And group "tea-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    When the user "Brian" gets all users of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |

  @skipOnStable2.0
  Scenario: admin user gets all users with certain roles and members of a certain group
    Given the administrator has given "Alice" the role "Admin" using the settings api
    And user "Carol" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "Space Admin" using the settings api
    And the administrator has given "Carol" the role "Space Admin" using the settings api
    And group "tea-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    When the user "Alice" gets all users with role "Space Admin" using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain following users with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |
      | Carol King   | %uuid_v4% | carol@example.org | Carol                    | true           |
    But the API response should not contain following user with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName |
      | Alice Hansen | %uuid_v4% | alice@example.org | Alice                    |
    When the user "Alice" gets all users with role "Space Admin" and member of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain following users with the information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    | true           |
    But the API response should not contain following user with the information:
      | displayName | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Carol King  | %uuid_v4% | carol@example.org | Carol                    | true           |

  @skipOnStable2.0
  Scenario Outline: non-admin user tries to get users with a certain role
    Given the administrator has given "Alice" the role "<userRole>" using the settings api
    When the user "Alice" gets all users with role "<role>" using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response
    Examples:
      | userRole    | role        |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | Guest       |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | Guest       |
      | User        | Admin       |
      | Guest       | Space Admin |
      | Guest       | User        |
      | Guest       | Guest       |
      | Guest       | Admin       |
