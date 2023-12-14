Feature: checking file versions using file id
  As a user
  I want to share file outside of the space
  So that other users can access the file


  Scenario Outline: check the file versions of a shared file
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "Project1" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project1" with content "hello world version 1" to "text.txt"
    And we save it into "FILEID"
    And user "Alice" has uploaded a file inside space "Project1" with content "hello world version 1.1" to "text.txt"
    And user "Alice" has created a share inside of space "Project1" with settings:
      | path      | text.txt |
      | shareWith | Brian    |
      | role      | <role>   |
    And using new DAV path
    When user "Alice" gets the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    When user "Brian" tries to get the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "403"
    Examples:
      | role   | 
      | editor | 
      | viewer |
      | all    |
