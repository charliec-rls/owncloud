@api
Feature: Propfind test
  As a user
  I want to check the PROPFIND response
  So that I can make sure that the response contains all the relevant values

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the GraphApi


  Scenario: space-admin checks the PROPFIND request of a space
    Given user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    When user "Alice" sends PROPFIND request to space "new-space" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "new-space" with these key and value pairs:
      | key            | value            |
      | oc:fileid      | UUIDof:new-space |
      | oc:name        | new-space        |
      | oc:permissions | RDNVCKZ          |
      | oc:privatelink |                  |
      | oc:size        | 12               |


  Scenario Outline: space member with a different role checks the PROPFIND request of a space
    Given user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    And user "Alice" has shared a space "new-space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" sends PROPFIND request to space "new-space" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "new-space" with these key and value pairs:
      | key            | value            |
      | oc:fileid      | UUIDof:new-space |
      | oc:name        | new-space        |
      | oc:permissions | <oc_permission>  |
      | oc:privatelink |                  |
      | oc:size        | 12               |
    Examples:
      | role    | oc_permission |
      | manager | RDNVCKZ       |
      | editor  | DNVCK         |
      | viewer  |               |
