@all @create_note
Feature: Create note
  In order to create a note
  The message should be passed as a request body parameter with minimum x characters
  The note must be created in the database
  I need to receive the created note back with the id and the message
  I need to be able to see the payer persisted on the DB with the correct message

  Background:
    Given the header is empty

  @test
  Scenario Outline: Create note error when <scenario> - response fields validation
    Given the "user" exists
    """
    {
      "id": "90d78048-f39d-47ab-9d24-3da4d8d1fb23",
      "name": "John Doe"
    }
    """
    When I call "POST" "/v1/users/90d78048-f39d-47ab-9d24-3da4d8d1fb23/notes" with body
    """
    {
      "title": <title>,
      "content": <content>
    }
    """
    Then the status returned should be 400
    And the response should contain the field "status_code" equal to "400"
    And the response should contain the field "messages.0" equal to "<error>"
    And the response should contain the field "type" equal to "Validation Error"

    Examples:
      | scenario                           | title | content | error                                 |
      | title is null                      | null  | "abc"   | 'title' is required                   |
      | title has less than 3 characters   | "ab"  | null    | 'title' should be greater in length   |
      | content is null                    | "abc" | null    | 'content' is required                 |
      | content has less than 3 characters | "abc" | "ab"    | 'content' should be greater in length |

  Scenario: Create note success - response fields validation
    Given the "user" exists
    """
    {
      "id": "90d78048-f39d-47ab-9d24-3da4d8d1fb23",
      "name": "John Doe"
    }
    """
    When I call "POST" "/v1/users/90d78048-f39d-47ab-9d24-3da4d8d1fb23/notes" with body
    """
    {
       "title": "My first note",
       "content": "This is my first note"
    }
    """
    Then the status returned should be 201
    And the response should contain the field "id" equal to "not nil"
    And the response should contain the field "title" equal to "My first note"
    And the response should contain the field "content" equal to "This is my first note"

  Scenario: Create note success - db fields validation
    Given the db should contain 0 objects in the "notes" table
    And the "user" exists
    """
    {
      "id": "90d78048-f39d-47ab-9d24-3da4d8d1fb23",
      "name": "John Doe"
    }
    """
    When I call "POST" "/v1/users/90d78048-f39d-47ab-9d24-3da4d8d1fb23/notes" with body
    """
    {
       "title": "My first note",
       "content": "This is my first note"
    }
    """
    Then the status returned should be 201
    And the db should contain 1 objects in the "notes" table
    And the db should contain the "note" with the "title" column value "My first note" colum "content" equal to "This is my first note"

  Scenario: Create note success - sqs fields validation
    Given the db should contain 0 objects in the "notes" table
    And the "user" exists
    """
    {
      "id": "90d78048-f39d-47ab-9d24-3da4d8d1fb23",
      "name": "John Doe"
    }
    """
    When I call "POST" "/v1/users/90d78048-f39d-47ab-9d24-3da4d8d1fb23/notes" with body
    """
    {
       "title": "My first note",
       "content": "This is my first note"
    }
    """
    Then the status returned should be 201
    And the sqs queue "https://localhost.localstack.cloud:4566/000000000000/newNoteQueue" should have 1 messages published
    And the sqs queue "https://localhost.localstack.cloud:4566/000000000000/newNoteQueue" should have a message published with "title" field equal to "My first note"
    And the sqs queue "https://localhost.localstack.cloud:4566/000000000000/newNoteQueue" should have a message published with "content" field equal to "This is my first note"