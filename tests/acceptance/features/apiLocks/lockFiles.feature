Feature: lock files
  As a user
  I want to lock files

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: lock a file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    When user "Alice" locks file "textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Infinity     |
      | d:lockdiscovery/d:activelock/oc:ownername            | Alice Hansen |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file with a timeout
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    When user "Alice" locks file "textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-5000 |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-5000  |
      | d:lockdiscovery/d:activelock/oc:ownername            | Alice Hansen |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file using file-id
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" locks file "textfile.txt" using file-id path "<dav-path>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-3600  |
      | d:lockdiscovery/d:activelock/oc:ownername            | Alice Hansen |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: user cannot lock file twice
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And user "Alice" has locked file "textfile.txt" setting the following properties
      | lockscope | exclusive |
    When user "Alice" tries to lock file "textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "423"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file in the project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "Project" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "textfile.txt"
    And user "Alice" has shared a space "Project" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" locks file "textfile.txt" inside the space "Project" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Brian" sends PROPFIND request from the space "Project" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Project" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-3600  |
      | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |
    Examples:
      | role    |
      | manager |
      | editor  |


  Scenario Outline: lock a file in the project space using file-id
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "Project" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "Project" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" locks file "textfile.txt" using file-id path "<dav-path>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Brian" sends PROPFIND request from the space "Project" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Project" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-3600  |
      | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |
    Examples:
      | role    | dav-path                          |
      | manager | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | /dav/spaces/<<FILEID>>            |


  Scenario: viewer cannot lock a file in the project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "Project" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "Project" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    When user "Brian" tries to lock file "textfile.txt" using file-id path "/dav/spaces/<<FILEID>>" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "403"
    When user "Brian" tries to lock file "textfile.txt" inside the space "Project" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "403"


  Scenario Outline: lock a file in the shares
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | editor       |
    When user "Brian" locks file "/Shares/textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file in the shares using file-id
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | editor       |
    When user "Brian" locks file "textfile.txt" using file-id path "<dav-path>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario: viewer cannot lock a file in the shares using file-id
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | viewer       |
    When user "Brian" tries to lock file "textfile.txt" using file-id path "/dav/spaces/<<FILEID>>" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "403"


  Scenario: sharee cannot lock a resource exclusively locked by a sharer
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | editor       |
    And user "Alice" has locked file "textfile.txt" setting the following properties
      | lockscope | exclusive |
    When user "Brian" tries to lock file "textfile.txt" using file-id path "/dav/spaces/<<FILEID>>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "423"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/oc:ownername            | Alice Hansen |


  Scenario: sharer cannot lock a resource exclusively locked by a sharee
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | editor       |
    And user "Brian" has locked file "textfile.txt" using file-id path "/dav/spaces/<<FILEID>>" setting the following properties
      | lockscope | exclusive |
    When user "Alice" tries to lock file "textfile.txt" using file-id path "/dav/spaces/<<FILEID>>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "423"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |

  @issue-7599
  Scenario Outline: two users having both a shared lock can use the resource
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "some data" to "textfile0.txt"
    And user "Brian" has uploaded file with content "some data" to "textfile0.txt"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian"
    And user "Alice" has locked file "textfile0.txt" setting the following properties
      | lockscope | shared |
    And user "Brian" has locked file "Shares/textfile0.txt" setting the following properties
      | lockscope | shared |
    When user "Alice" uploads file with content "from user 0" to "textfile0.txt" sending the locktoken of file "textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And the content of file "textfile0.txt" for user "Alice" should be "from user 0"
    And the content of file "Shares/textfile0.txt" for user "Brian" should be "from user 0"
    When user "Brian" uploads file with content "from user 1" to "Shares/textfile0.txt" sending the locktoken of file "Shares/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And the content of file "textfile0.txt" for user "Alice" should be "from user 1"
    And the content of file "Shares/textfile0.txt" for user "Brian" should be "from user 1"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-7638
  Scenario Outline: locking a file does not lock other items with the same name in other parts of the file system
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "locked"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/locked/textfile0.txt"
    And user "Alice" has created folder "notlocked"
    And user "Alice" has created folder "notlocked/textfile0.txt"
    When user "Alice" locks file "locked/textfile0.txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "200"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/notlocked/textfile0.txt/real-file.txt"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/textfile0.txt"
    But user "Alice" should not be able to upload file "filesForUpload/lorem.txt" to "/locked/textfile0.txt"
    Examples:
      | dav-path-version | lock-scope |
      | old              | shared     |
      | old              | exclusive  |
      | new              | shared     |
      | new              | exclusive  |
      | spaces           | shared     |
      | spaces           | exclusive  |

  @issue-7638 @issue-7599
  Scenario Outline: locking a file in a received share does not lock other items with the same name in other received shares (shares from different users)
    Given using <dav-path-version> DAV path
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "FromAlice"
    And user "Brian" has created folder "FromBrian"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/FromAlice/textfile0.txt"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "/FromBrian/textfile0.txt"
    And user "Alice" has shared folder "FromAlice" with user "Carol"
    And user "Brian" has shared folder "FromBrian" with user "Carol"
    When user "Carol" locks file "/Shares/FromBrian/textfile0.txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "200"
    And user "Carol" should be able to upload file "filesForUpload/lorem.txt" to "/Shares/FromAlice/textfile0.txt"
    But user "Carol" should not be able to upload file "filesForUpload/lorem.txt" to "/Shares/FromBrian/textfile0.txt"
    Examples:
      | dav-path-version | lock-scope |
      | old              | shared     |
      | old              | exclusive  |
      | new              | shared     |
      | new              | exclusive  |
      | spaces           | shared     |
      | spaces           | exclusive  |

  @issue-7638 @issue-7599
  Scenario Outline: locking a file in a received share does not lock other items with the same name in other received shares (shares from same user)
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "locked/"
    And user "Alice" has created folder "notlocked/"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/locked/textfile0.txt"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/notlocked/textfile0.txt"
    And user "Alice" has shared folder "locked" with user "Brian"
    And user "Alice" has shared folder "notlocked" with user "Brian"
    When user "Brian" locks file "/Shares/locked/textfile0.txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "200"
    And user "Brian" should be able to upload file "filesForUpload/lorem.txt" to "/Shares/notlocked/textfile0.txt"
    But user "Brian" should not be able to upload file "filesForUpload/lorem.txt" to "/Shares/locked/textfile0.txt"
    Examples:
      | dav-path-version | lock-scope |
      | old              | shared     |
      | old              | exclusive  |
      | new              | shared     |
      | new              | exclusive  |
      | spaces           | shared     |
      | spaces           | exclusive  |

  @issue-7641
  Scenario Outline: try to lock a folder
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "locked/"
    When user "Alice" tries to lock folder "locked/" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "403"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "locked/textfile0.txt"
    And user "Alice" should be able to create folder "/locked/sub-folder"
    Examples:
      | dav-path-version | lock-scope |
      | old              | shared     |
      | old              | exclusive  |
      | new              | shared     |
      | new              | exclusive  |
      | spaces           | shared     |
      | spaces           | exclusive  |
