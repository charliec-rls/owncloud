@api
Feature: create user
  As a admin
  I want to create a user
  So that the user can use the application

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-3516
  Scenario Outline: admin creates a user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | <userName>    |
      | displayName    | <displayName> |
      | email          | <email>       |
      | password       | <password>    |
      | accountEnabled | <enable>      |
    Then the HTTP status code should be "<code>"
    And user "<userName>" <shouldOrNot> exist
    Examples:
      | userName                     | displayName     | email                   | password                     | code | enable | shouldOrNot |
      | SameDisplayName              | Alice Hansen    | new@example.org         | containsCharacters(*:!;_+-&) | 200  | true   | should      |
      | withoutPassSameEmail         | without pass    | alice@example.org       |                              | 200  | true   | should      |
      | name                         | pass with space | example@example.org     | my pass                      | 200  | true   | should      |
      | nameWithCharacters(*:!;_+-&) | user            | new@example.org         | 123                          | 400  | true   | should not  |
      | name with space              | name with space | example@example.org     | 123                          | 400  | true   | should not  |
      | createDisabledUser           | disabled user   | example@example.org     | 123                          | 200  | false  | should      |
      | nameWithNumbers0123456       | user            | name0123456@example.org | 123                          | 200  | true   | should      |
      | name.with.dots               | user            | name.w.dots@example.org | 123                          | 200  | true   | should      |
      | 123456789                    | user            | 123456789@example.org   | 123                          | 400  | true   | should not  |
      | 0.0                          | user            | float@example.org       | 123                          | 400  | true   | should not  |

    @skipOnStable2.0
    Examples:
      | userName                     | displayName     | email               | password                     | code | enable | shouldOrNot |
      | withoutEmail                 | without email   |                     | 123                          | 200  | true   | should      |
      | Alice                        | same userName   | new@example.org     | 123                          | 409  | true   | should      |


  Scenario: user cannot be created with empty name
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       |              |
      | displayName    | emptyName    |
      | email          | @example.org |
      | password       | 123          |
      | accountEnabled | true         |
    Then the HTTP status code should be "400"


  Scenario Outline: user without admin right cannot create a user
    Given the administrator has assigned the role "<role>" to user "Alice" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | user         |
      | displayName    | user         |
      | email          | @example.org |
      | password       | 123          |
      | accountEnabled | true         |
    Then the HTTP status code should be "401"
    And user "user" should not exist
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: user cannot be created with the name of the disabled user
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the user "Alice" has disabled user "Brian" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | Brian                 |
      | displayName    | This is another Brian |
      | email          | brian@example.com     |
      | password       | 123                   |
      | accountEnabled | true                  |
    Then the HTTP status code should be "409"

  @skipOnStable2.0
  Scenario: user can be created with the name of the deleted user
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the user "Alice" has deleted a user "Brian" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | Brian                 |
      | displayName    | This is another Brian |
      | email          | brian@example.com     |
      | password       | 123                   |
      | accountEnabled | true                  |
    Then the HTTP status code should be "200"
    And user "Brian" should exist
