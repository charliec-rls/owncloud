Feature: download file
  As a user
  I want to be able to download files
  So that I can work wih local copies of files on my client system

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "Welcome this is just an example file for developers." to "/welcome.txt"

  @smokeTest
  Scenario Outline: download a file
    Given using <dav-path-version> DAV path
    When user "Alice" downloads file "/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded content should be "ownCloud test text file 0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @issue-1346
  Scenario Outline: download a file with range
    Given using <dav-path-version> DAV path
    When user "Alice" downloads file "/welcome.txt" with range "bytes=24-50" using the WebDAV API
    Then the HTTP status code should be "206"
    And the downloaded content should be "example file for developers"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: download a file larger than 4MB (ref: https://github.com/sabre-io/http/pull/119 )
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "/file9000000.txt" ending with "text at end of file" of size 9000000 bytes
    When user "Alice" downloads file "/file9000000.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the size of the downloaded file should be 9000000 bytes
    And the downloaded content should end with "text at end of file"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: get the size of a file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "This is a test file" to "test-file.txt"
    When user "Alice" gets the size of file "test-file.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the size of the file should be "19"

    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @issue-1316
  Scenario Outline: get the content-length response header of a pdf file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/simple.pdf" to "/simple.pdf"
    When user "Alice" downloads file "/simple.pdf" using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header         | value |
      | Content-Length | 9622  |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @issue-1316
  Scenario Outline: get the content-length response header of an image file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "/testavatar.png"
    When user "Alice" downloads file "/testavatar.png" using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header         | value |
      | Content-Length | 35323 |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: download a file with comma in the filename
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "file with comma in filename" to <file-name>
    When user "Alice" downloads file <file-name> using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded content should be "file with comma in filename"
    Examples:
      | dav-path-version | file-name      |
      | old              | "sample,1.txt" |
      | old              | ",,,.txt"      |
      | old              | ",,,.,"        |
      | new              | "sample,1.txt" |
      | new              | ",,,.txt"      |
      | new              | ",,,.,"        |

    @skipOnRevaMaster
    Examples:
      | dav-path-version | file-name      |
      | spaces           | "sample,1.txt" |
      | spaces           | ",,,.txt"      |
      | spaces           | ",,,.,"        |


  Scenario Outline: download a file with single part ranges
    Given using <dav-path-version> DAV path
    When user "Alice" downloads file "/welcome.txt" with range "bytes=0-51" using the WebDAV API
    Then the HTTP status code should be "206"
    And the following headers should be set
      | header         | value         |
      | Content-Length | 52            |
      | Content-Range  | bytes 0-51/52 |
    And the downloaded content should be "Welcome this is just an example file for developers."
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: download a file with multipart ranges
    Given using <dav-path-version> DAV path
    When user "Alice" downloads file "/welcome.txt" with range "bytes=0-6, 40-51" using the WebDAV API
    Then the HTTP status code should be "206" or "200"
    And if the HTTP status code was "206" then the following headers should match these regular expressions
      | Content-Length | /\d+/                                               |
      | Content-Type   | /^multipart\/byteranges; boundary=[a-zA-Z0-9_.-]*$/ |
    And if the HTTP status code was "206" then the downloaded content for multipart byterange should be:
      """
      Content-type: text/plain;charset=UTF-8
      Content-range: bytes 0-6/52

      Welcome

      Content-type: text/plain;charset=UTF-8
      Content-range: bytes 40-51/52

      developers.
      """
    But if the HTTP status code was "200" then the downloaded content should be "Welcome this is just an example file for developers."
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: download a file with last byte range out of bounds
    Given using <dav-path-version> DAV path
    When user "Alice" downloads file "/welcome.txt" with range "bytes=0-55" using the WebDAV API
    Then the HTTP status code should be "206"
    And the downloaded content should be "Welcome this is just an example file for developers."
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: download a range at the end of a file
    Given using <dav-path-version> DAV path
    When user "Alice" downloads file "/welcome.txt" with range "bytes=-11" using the WebDAV API
    Then the HTTP status code should be "206"
    And the downloaded content should be "developers."
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: download a file with range out of bounds
    Given using <dav-path-version> DAV path
    When user "Alice" downloads file "/welcome.txt" with range "bytes=55-60" using the WebDAV API
    Then the HTTP status code should be "416"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: download hidden files
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has uploaded the following files with content "hidden file"
      | path                 |
      | .hidden_file         |
      | /FOLDER/.hidden_file |
    When user "Alice" downloads the following files using the WebDAV API
      | path                 |
      | .hidden_file         |
      | /FOLDER/.hidden_file |
    Then the HTTP status code of responses on all endpoints should be "200"
    And the content of the following files for user "Alice" should be "hidden file"
      | path                 |
      | .hidden_file         |
      | /FOLDER/.hidden_file |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @smokeTest @issue-8361 @skipOnReva
  Scenario Outline: downloading a file should serve security headers
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "test file" to "/<file-name>"
    When user "Alice" downloads file "/<file-name>" using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header                            | value                                                              |
      | Content-Disposition               | attachment; filename*=UTF-8''"<file-name>"; filename="<file-name>" |
      | Content-Security-Policy           | child-src 'self'; connect-src 'self'; default-src 'none'; font-src 'self'; frame-ancestors 'self'; frame-src 'self' blob: https://embed.diagrams.net/; img-src 'self' data: blob:; manifest-src 'self'; media-src 'self'; object-src 'self' blob:; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'                                                |
      | X-Content-Type-Options            | nosniff                                                            |
      | X-Download-Options                | noopen                                                             |
      | X-Frame-Options                   | SAMEORIGIN                                                         |
      | X-Permitted-Cross-Domain-Policies | none                                                               |
      | X-Robots-Tag                      | none                                                               |
      | X-XSS-Protection                  | 1; mode=block                                                      |
    And the downloaded content should be "test file"
    Examples:
      | dav-path-version | file-name          |
      | old              | textfile.txt       |
      | old              | comma,.txt         |
      | old              | 'quote'single'.txt |
      | new              | textfile.txt       |
      | new              | comma,.txt         |
      | new              | 'quote'single'.txt |
      | spaces           | textfile.txt       |
      | spaces           | comma,.txt         |
      | spaces           | 'quote'single'.txt |

  @smokeTest @issue-8361 @skipOnReva
  Scenario Outline: downloading a file should serve security headers (file with doubel quotes)
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "test file" to '/"quote"double".txt'
    When user "Alice" downloads file '/"quote"double".txt' using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header                            | value                                                                            |
      | Content-Disposition               | attachment; filename*=UTF-8''""quote"double".txt"; filename=""quote"double".txt" |
      | Content-Security-Policy           | child-src 'self'; connect-src 'self'; default-src 'none'; font-src 'self'; frame-ancestors 'self'; frame-src 'self' blob: https://embed.diagrams.net/; img-src 'self' data: blob:; manifest-src 'self'; media-src 'self'; object-src 'self' blob:; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'                                                              |
      | X-Content-Type-Options            | nosniff                                                                          |
      | X-Download-Options                | noopen                                                                           |
      | X-Frame-Options                   | SAMEORIGIN                                                                       |
      | X-Permitted-Cross-Domain-Policies | none                                                                             |
      | X-Robots-Tag                      | none                                                                             |
      | X-XSS-Protection                  | 1; mode=block                                                                    |
    And the downloaded content should be "test file"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario: download a zero byte size file
    Given user "Alice" has uploaded file "filesForUpload/zerobyte.txt" to "/zerobyte.txt"
    When user "Alice" downloads file "/zerobyte.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the size of the downloaded file should be 0 bytes


  Scenario: try to download recently deleted file
    When user "Alice" deletes file "textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    When user "Alice" tries to download file "textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "404"
