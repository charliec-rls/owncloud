@api @skipOnOcV10
Feature: create file or folder named similar to Shares folder
  As a user
  I want to be able to create files and folders when the Shares folder exists
  So that I can organise the files in my file system

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has shared folder "/FOLDER" with user "Brian" with permissions "read,update"
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"

  Scenario Outline: create a folder with a name similar to Shares
    Given using spaces DAV path
    When user "Brian" creates folder "<folder_name>" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Brian" the space "Personal" should contain these entries:
      | <folder_name>/ |
    Examples:
      | folder_name |
      | Share       |
      | shares      |
      | Share1      |

  Scenario Outline: create a file with a name similar to Shares
    Given using spaces DAV path
    When user "Brian" uploads file with content "some text" to "<file_name>" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "<file_name>" for user "Brian" should be "some text"
    And for user "Brian" the space "Personal" should contain these entries:
      | <file_name> |
    And for user "Brian" the space "Shares Jail" should contain these entries:
      | FOLDER/ |
    Examples:
      | file_name |
      | Share    |
      | shares   |
      | Share1   |

  Scenario: try to create a folder named Shares
    Given using spaces DAV path
    When user "Brian" creates folder "/Shares" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Brian" the space "Shares Jail" should contain these entries:
      | FOLDER/ |

  Scenario: try to create a file named Shares
    Given using spaces DAV path
    When user "Brian" uploads file with content "some text" to "/Shares" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Brian" the space "Shares Jail" should contain these entries:
      | FOLDER/ |
