@api @skipOnOcV10
Feature: copy file
  As a user
  I want to be able to copy files
  So that I can manage my files

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: Copying a file within a same space project with role manager and editor
    Given the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project" with the default quota using the GraphApi
    And user "Alice" has created a folder "newfolder" in space "Project"
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "insideSpace.txt"
    And user "Alice" has shared a space "Project" to user "Brian" with role "<role>"
    When user "Brian" copies file "insideSpace.txt" to "/newfolder/insideSpace.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Brian" folder "newfolder" of the space "Project" should contain these entries:
      | insideSpace.txt |
    And for user "Alice" the content of the file "newfolder/insideSpace.txt" of the space "Project" should be "some content"
    Examples:
      | role    |
      | manager |
      | editor  |


  Scenario: Copying a file within a same space project with role viewer
    Given the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project" with the default quota using the GraphApi
    And user "Alice" has created a folder "newfolder" in space "Project"
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "insideSpace.txt"
    And user "Alice" has shared a space "Project" to user "Brian" with role "viewer"
    When user "Brian" copies file "insideSpace.txt" to "newfolder/insideSpace.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Brian" the space "Project" should not contain these entries:
      | newfolder/insideSpace.txt |


  Scenario Outline: User copies a file from a space project with a different role to a space project with the manager role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project1" with the default quota using the GraphApi
    And user "Brian" has created a space "Project2" with the default quota using the GraphApi
    And user "Brian" has uploaded a file inside space "Project1" with content "Project1 content" to "project1.txt"
    And user "Brian" has shared a space "Project2" to user "Alice" with role "<to_role>"
    And user "Brian" has shared a space "Project1" to user "Alice" with role "<from_role>"
    When user "Alice" copies file "project1.txt" from space "Project1" to "project1.txt" inside space "Project2" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project2" should contain these entries:
      | project1.txt |
    And for user "Alice" the content of the file "project1.txt" of the space "Project2" should be "Project1 content"
    Examples:
      | from_role | to_role |
      | manager   | manager |
      | manager   | editor  |
      | editor    | manager |
      | editor    | editor  |


  Scenario Outline: User copies a file from a space project with a different role to a space project with a viewer role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project1" with the default quota using the GraphApi
    And user "Brian" has created a space "Project2" with the default quota using the GraphApi
    And user "Brian" has uploaded a file inside space "Project1" with content "Project1 content" to "project1.txt"
    And user "Brian" has shared a space "Project2" to user "Alice" with role "viewer"
    And user "Brian" has shared a space "Project1" to user "Alice" with role "<role>"
    When user "Alice" copies file "project1.txt" from space "Project1" to "project1.txt" inside space "Project2" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Project2" should not contain these entries:
      | project1.txt |
    Examples:
      | role    |
      | manager |
      | editor  |


  Scenario Outline: User copies a file from space project with different role to space personal
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has uploaded a file inside space "Project" with content "Project content" to "project.txt"
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    When user "Alice" copies file "project.txt" from space "Project" to "project.txt" inside space "Personal" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | project.txt |
    And for user "Alice" the content of the file "project.txt" of the space "Personal" should be "Project content"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: User copies a file from space project with different role to space shares jail with editor role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded a file inside space "Project" with content "Project content" to "project.txt"
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "31"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" copies file "project.txt" from space "Project" to "/testshare/project.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Shares Jail" should contain these entries:
      | /testshare/project.txt |
    And for user "Alice" the content of the file "/testshare/project.txt" of the space "Shares Jail" should be "Project content"
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: User copies a file from space project with different role to shares jail with viewer role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded a file inside space "Project" with content "Project content" to "project.txt"
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "17"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" copies file "project.txt" from space "Project" to "/testshare/project.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Shares Jail" should not contain these entries:
      | /testshare/project.txt |
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario Outline: User copies a file from space personal to space project with different role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    And user "Alice" has uploaded file with content "personal space content" to "/personal.txt"
    When user "Alice" copies file "personal.txt" from space "Personal" to "personal.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project" should contain these entries:
      | personal.txt |
    And for user "Alice" the content of the file "personal.txt" of the space "Project" should be "personal space content"
    Examples:
      | role    |
      | manager |
      | editor  |


  Scenario: User copies a file from space personal to space project with role viewer
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has shared a space "Project" to user "Alice" with role "viewer"
    And user "Alice" has uploaded file with content "personal space content" to "/personal.txt"
    When user "Alice" copies file "personal.txt" from space "Personal" to "personal.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Project" should not contain these entries:
      | personal.txt |


  Scenario: User copies a file from space personal to space shares jail with role editor
    Given user "Brian" has created folder "/testshare"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "31"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    And user "Alice" has uploaded file with content "personal content" to "personal.txt"
    When user "Alice" copies file "personal.txt" from space "Personal" to "/testshare/personal.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Shares Jail" should contain these entries:
      | /testshare/personal.txt |
    And for user "Alice" the content of the file "/testshare/personal.txt" of the space "Shares Jail" should be "personal content"


  Scenario: User copies a file from space personal to space shares jail with role viewer
    Given user "Brian" has created folder "/testshare"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "17"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    And user "Alice" has uploaded file with content "personal content" to "personal.txt"
    When user "Alice" copies file "personal.txt" from space "Personal" to "/testshare/personal.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Shares Jail" should not contain these entries:
      | /testshare/personal.txt |


  Scenario Outline: User copies a file from space shares jail with different role to space personal
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "testshare content" to "/testshare/testshare.txt"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "<permissions>"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" copies file "/testshare/testshare.txt" from space "Shares Jail" to "testshare.txt" inside space "Personal" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | /testshare.txt |
    And for user "Alice" the content of the file "/testshare.txt" of the space "Personal" should be "testshare content"
    Examples:
      | permissions |
      | 31          |
      | 17          |


  Scenario Outline: User copies a file from space shares jail with different role to space project with different role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "testshare content" to "/testshare/testshare.txt"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "<permissions>"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" copies file "/testshare/testshare.txt" from space "Shares Jail" to "testshare.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project" should contain these entries:
      | /testshare.txt |
    And for user "Alice" the content of the file "testshare.txt" of the space "Project" should be "testshare content"
    Examples:
      | role    | permissions |
      | manager | 31          |
      | manager | 17          |
      | editor  | 31          |
      | editor  | 17          |


  Scenario Outline: User copies a file from space shares jail with different role to space project with role viewer
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has shared a space "Project" to user "Alice" with role "viewer"
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "testshare content" to "/testshare/testshare.txt"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "<permissions>"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" copies file "/testshare/testshare.txt" from space "Shares Jail" to "testshare.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Project" should not contain these entries:
      | /testshare.txt |
    Examples:
      | permissions |
      | 31          |
      | 17          |


  Scenario Outline: User copies a file from space shares jail with different role to space shares jail with role editor
    Given user "Brian" has created folder "/testshare1"
    And user "Brian" has created folder "/testshare2"
    And user "Brian" has uploaded file with content "testshare1 content" to "/testshare1/testshare1.txt"
    And user "Brian" has shared folder "/testshare1" with user "Alice" with permissions "<permissions>"
    And user "Brian" has shared folder "/testshare2" with user "Alice" with permissions "31"
    And user "Alice" has accepted share "/testshare1" offered by user "Brian"
    And user "Alice" has accepted share "/testshare2" offered by user "Brian"
    When user "Alice" copies file "/testshare1/testshare1.txt" from space "Shares Jail" to "/testshare2/testshare1.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Shares Jail" should contain these entries:
      | /testshare2/testshare1.txt |
    And for user "Brian" the space "Personal" should contain these entries:
      | /testshare2/testshare1.txt |
    And for user "Alice" the content of the file "/testshare2/testshare1.txt" of the space "Shares Jail" should be "testshare1 content"
    And for user "Brian" the content of the file "/testshare1/testshare1.txt" of the space "Personal" should be "testshare1 content"
    Examples:
      | permissions |
      | 31          |
      | 17          |


  Scenario Outline: User copies a file from space shares jail with different role to space shares jail with role editor
    Given user "Brian" has created folder "/testshare1"
    And user "Brian" has created folder "/testshare2"
    And user "Brian" has uploaded file with content "testshare1 content" to "/testshare1/testshare1.txt"
    And user "Brian" has shared folder "/testshare1" with user "Alice" with permissions "<permissions>"
    And user "Brian" has shared folder "/testshare2" with user "Alice" with permissions "17"
    And user "Alice" has accepted share "/testshare1" offered by user "Brian"
    And user "Alice" has accepted share "/testshare2" offered by user "Brian"
    When user "Alice" copies file "/testshare1/testshare1.txt" from space "Shares Jail" to "/testshare2/testshare1.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Shares Jail" should not contain these entries:
      | /testshare2/testshare1.txt |
    And for user "Brian" the space "Personal" should not contain these entries:
      | /testshare2/testshare1.txt |
    Examples:
      | permissions |
      | 31          |
      | 17          |
