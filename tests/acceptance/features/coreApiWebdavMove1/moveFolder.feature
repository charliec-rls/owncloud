@api
Feature: move (rename) folder
  As a user
  I want to be able to move and rename folders
  So that I can quickly manage my file system

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: renaming a folder to a backslash should return an error
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "\" using the WebDAV API
    Then the HTTP status code should be "400"
    And user "Alice" should see the following elements
      | /testshare/ |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: renaming a folder beginning with a backslash should return an error
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "\testshare" using the WebDAV API
    Then the HTTP status code should be "400"
    And user "Alice" should see the following elements
      | /testshare/ |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: renaming a folder including a backslash encoded should return an error
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "/hola\hola" using the WebDAV API
    Then the HTTP status code should be "400"
    And user "Alice" should see the following elements
      | /testshare/ |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: move a folder into an other folder
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/testshare"
    And user "Alice" has created folder "/an-other-folder"
    When user "Alice" moves folder "/testshare" to "/an-other-folder/testshare" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Alice" should not see the following elements
      | /testshare/ |
    And user "Alice" should see the following elements
      | /an-other-folder/testshare/ |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: move a folder into a nonexistent folder
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "/not-existing/testshare" using the WebDAV API
    Then the HTTP status code should be "409"
    And user "Alice" should see the following elements
      | /testshare/ |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: renaming folder with dots in the path
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "<folder_name>"
    And user "Alice" has uploaded file with content "uploaded content for file name ending with a dot" to "<folder_name>/abc.txt"
    When user "Alice" moves folder "<folder_name>" to "/uploadFolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/uploadFolder/abc.txt" for user "Alice" should be "uploaded content for file name ending with a dot"
    Examples:
      | dav-path-version | folder_name   |
      | old              | /upload.      |
      | old              | /upload.1     |
      | old              | /upload...1.. |
      | old              | /...          |
      | old              | /..upload     |
      | new              | /upload.      |
      | new              | /upload.1     |
      | new              | /upload...1.. |
      | new              | /...          |
      | new              | /..upload     |

    @skipOnRevaMaster
    Examples:
      | dav-path-version | folder_name   |
      | spaces           | /upload.      |
      | spaces           | /upload.1     |
      | spaces           | /upload...1.. |
      | spaces           | /...          |
      | spaces           | /..upload     |

  @skipOnRevaMaster @issue-3023
  Scenario Outline: moving a folder into a sub-folder of itself
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "PARENT"
    And user "Alice" has created folder "PARENT/CHILD"
    And user "Alice" has uploaded file with content "parent text" to "/PARENT/parent.txt"
    And user "Alice" has uploaded file with content "child text" to "/PARENT/CHILD/child.txt"
    When user "Alice" moves folder "/PARENT" to "/PARENT/CHILD/PARENT" using the WebDAV API
    Then the HTTP status code should be "409"
    And the content of file "/PARENT/parent.txt" for user "Alice" should be "parent text"
    And the content of file "/PARENT/CHILD/child.txt" for user "Alice" should be "child text"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
