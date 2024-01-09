@skipOnReva
Feature: move (rename) file
  As a user
  I want to be able to move and rename files
  So that I can manage my file system

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: sharer moves a file into a shared folder
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    And user "Brian" has uploaded file with content "test data" to "/testfile.txt"
    When user "Brian" moves file "/testfile.txt" to "testshare/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/Shares/testshare/testfile.txt" for user "Alice" should be "test data"
    And the content of file "/testshare/testfile.txt" for user "Brian" should be "test data"
    And as "Brian" file "/testfile.txt" should not exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: sharee tries to move a file into a shared folder
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    And user "Alice" has uploaded file with content "test data" to "/testfile.txt"
    When user "Alice" moves file "/testfile.txt" to "Shares/testshare/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "Shares/testshare/testfile.txt" should not exist
    And as "Brian" file "testshare/testfile.txt" should not exist
    But as "Alice" file "/testfile.txt" should exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: moving a file out of a shared folder as the sharer
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "test data" to "/testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    When user "Brian" moves file "/testshare/testfile.txt" to "/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/testfile.txt" for user "Brian" should be "test data"
    And as "Alice" file "/Shares/testshare/testfile.txt" should not exist
    And as "Brian" file "/testshare/testfile.txt" should not exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: moving a file out of a shared folder as the sharee
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "test data" to "/testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    When user "Alice" moves file "/Shares/testshare/testfile.txt" to "/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/Shares/testshare/testfile.txt" should exist
    And as "Brian" file "/testshare/testfile.txt" should exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: moving a folder into a shared folder the sharer
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    And user "Brian" has created folder "/testsubfolder"
    And user "Brian" has uploaded file with content "test data" to "/testsubfolder/testfile.txt"
    When user "Brian" moves folder "/testsubfolder" to "testshare/testsubfolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/Shares/testshare/testsubfolder/testfile.txt" for user "Alice" should be "test data"
    And the content of file "/testshare/testsubfolder/testfile.txt" for user "Brian" should be "test data"
    And as "Brian" file "/testsubfolder" should not exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: moving a folder into a shared folder as the sharee
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    And user "Alice" has created folder "/testsubfolder"
    And user "Alice" has uploaded file with content "test data" to "/testsubfolder/testfile.txt"
    When user "Alice" moves folder "/testsubfolder" to "Shares/testshare/testsubfolder" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" folder "/Shares/testshare/testsubfolder" should not exist
    And as "Brian" folder "/testshare/testsubfolder" should not exist
    But as "Alice" folder "/testsubfolder" should exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: moving a folder out of a shared folder as the sharer
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created the following folders
      | path                     |
      | /testshare               |
      | /testshare/testsubfolder |
    And user "Brian" has uploaded file with content "test data" to "/testshare/testsubfolder/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    When user "Brian" moves folder "/testshare/testsubfolder" to "/testsubfolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/testsubfolder/testfile.txt" for user "Brian" should be "test data"
    And as "Alice" folder "/testshare/testsubfolder" should not exist
    And as "Brian" folder "/testshare/testsubfolder" should not exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: moving a folder out of a shared folder as the sharee
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created the following folders
      | path                     |
      | /testshare               |
      | /testshare/testsubfolder |
    And user "Brian" has uploaded file with content "test data" to "/testshare/testsubfolder/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    When user "Alice" moves folder "/Shares/testshare/testsubfolder" to "/testsubfolder" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" folder "/Shares/testshare/testsubfolder" should exist
    And as "Brian" folder "/testshare/testsubfolder" should exist
    Examples:
      | dav-path-version | permissions |
      | old              | read        |
      | old              | change      |
      | old              | all         |
      | new              | read        |
      | new              | change      |
      | new              | all         |


  Scenario Outline: sharee moves a file within a shared folder (change/all permissions)
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "testshare"
    And user "Brian" has created folder "testshare/child"
    And user "Brian" has uploaded file with content "test data" to "testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    When user "Alice" moves folder "Shares/testshare/testfile.txt" to "Shares/testshare/child/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/Shares/testshare/child/testfile.txt" should exist
    And as "Brian" file "/testshare/child/testfile.txt" should exist
    And as "Alice" file "/Shares/testshare/testfile.txt" should not exist
    And as "Brian" file "/testshare/testfile.txt" should not exist
    Examples:
      | dav-path-version | permissions |
      | old              | change      |
      | old              | all         |
      | new              | change      |
      | new              | all         |


  Scenario Outline: sharee tries to move a file within a shared folder (read permissions)
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "testshare"
    And user "Brian" has created folder "testshare/child"
    And user "Brian" has uploaded file with content "test data" to "testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | read          |
      | shareWith   | Alice         |
    When user "Alice" moves folder "Shares/testshare/testfile.txt" to "Shares/testshare/child/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/Shares/testshare/child/testfile.txt" should not exist
    And as "Brian" file "/testshare/child/testfile.txt" should not exist
    And as "Alice" file "/Shares/testshare/testfile.txt" should exist
    And as "Brian" file "/testshare/testfile.txt" should exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @issue-1976
  Scenario Outline: sharee tries to move a file into same shared folder with same name
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "testshare"
    And user "Brian" has uploaded file with content "test data" to "testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare     |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Alice         |
    When user "Alice" moves folder "Shares/testshare/testfile.txt" to "Shares/testshare/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And the content of file "Shares/testshare/testfile.txt" for user "Alice" should be "test data"
    And the content of file "testshare/testfile.txt" for user "Brian" should be "test data"
    Examples:
      | dav-path-version | permissions |
      | old              | all         |
      | old              | change      |
      | old              | read        |
      | new              | all         |
      | new              | change      |
      | new              | read        |
