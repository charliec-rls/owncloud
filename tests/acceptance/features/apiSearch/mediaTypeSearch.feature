Feature: media type search
  As a user
  I want to search files using media type
  So that I can find the files with specific media type

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario Outline: search for files using media type
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/lorem.txt"
    And user "Alice" has uploaded file "filesForUpload/simple.pdf" to "/simple.pdf"
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "/testavatar.jpg"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "/testavatar.png"
    And user "Alice" has uploaded file "filesForUpload/data.tar.gz" to "/data.tar.gz"
    When user "Alice" searches for "mediatype:<pattern>" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | <search-result> |
    Examples:
      | pattern | search-result   |
      | *text*  | /lorem.txt      |
      | *pdf*   | /simple.pdf     |
      | *jpeg*  | /testavatar.jpg |
      | *png*   | /testavatar.png |
      | *gzip*  | /data.tar.gz    |


  Scenario Outline: search for files inside sub folders using media type
    Given user "Alice" has created folder "/uploadFolder"
    And user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/uploadFolder/lorem.txt"
    And user "Alice" has uploaded file "filesForUpload/simple.pdf" to "/uploadFolder/simple.pdf"
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "/uploadFolder/testavatar.jpg"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "/uploadFolder/testavatar.png"
    And user "Alice" has uploaded file "filesForUpload/data.tar.gz" to "/uploadFolder/data.tar.gz"
    When user "Alice" searches for "mediatype:<pattern>" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | <search-result> |
    Examples:
      | pattern | search-result                |
      | *text*  | /uploadFolder/lorem.txt      |
      | *pdf*   | /uploadFolder/simple.pdf     |
      | *jpeg*  | /uploadFolder/testavatar.jpg |
      | *png*   | /uploadFolder/testavatar.png |
      | *gzip*  | /uploadFolder/data.tar.gz    |


  Scenario Outline: search for file inside project space using media type
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "find data" with the default quota using the Graph API
    And user "Alice" has uploaded a file "filesForUpload/lorem.txt" to "/lorem.txt" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/simple.pdf" to "/simple.pdf" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "/testavatar.jpg" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.png" to "/testavatar.png" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/data.tar.gz" to "/data.tar.gz" in space "find data"
    When user "Alice" searches for "mediatype:<pattern>" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | <search-result> |
    Examples:
      | pattern | search-result   |
      | *text*  | /lorem.txt      |
      | *pdf*   | /simple.pdf     |
      | *jpeg*  | /testavatar.jpg |
      | *png*   | /testavatar.png |
      | *gzip*  | /data.tar.gz    |


  Scenario Outline: sharee searches for shared files using media type
    Given user "Alice" has created folder "/uploadFolder"
    And user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/uploadFolder/lorem.txt"
    And user "Alice" has uploaded file "filesForUpload/simple.pdf" to "/uploadFolder/simple.pdf"
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "/uploadFolder/testavatar.jpg"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "/uploadFolder/testavatar.png"
    And user "Alice" has uploaded file "filesForUpload/data.tar.gz" to "/uploadFolder/data.tar.gz"
    And user "Alice" has shared folder "uploadFolder" with user "Brian"
    When user "Brian" searches for "mediatype:<pattern>" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | <search-result> |
    Examples:
      | pattern | search-result                |
      | *text*  | /uploadFolder/lorem.txt      |
      | *pdf*   | /uploadFolder/simple.pdf     |
      | *jpeg*  | /uploadFolder/testavatar.jpg |
      | *png*   | /uploadFolder/testavatar.png |
      | *gzip*  | /uploadFolder/data.tar.gz    |


  Scenario Outline: sharee searches for files inside shared space using media type
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "find data" with the default quota using the Graph API
    And user "Alice" has uploaded a file "filesForUpload/lorem.txt" to "/lorem.txt" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/simple.pdf" to "/simple.pdf" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "/testavatar.jpg" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.png" to "/testavatar.png" in space "find data"
    And user "Alice" has uploaded a file "filesForUpload/data.tar.gz" to "/data.tar.gz" in space "find data"
    And user "Alice" has shared a space "find data" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    When user "Brian" searches for "mediatype:<pattern>" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | <search-result> |
    Examples:
      | pattern | search-result   |
      | *text*  | /lorem.txt      |
      | *pdf*   | /simple.pdf     |
      | *jpeg*  | /testavatar.jpg |
      | *png*   | /testavatar.png |
      | *gzip*  | /data.tar.gz    |
