@all @get_note
Feature: Get note

  Background:
    Given the header is empty

  Scenario: Get note success - response fields validation
    Given the "user" exists
    """
    {
      "id": "90d78048-f39d-47ab-9d24-3da4d8d1fb23",
      "name": "John Doe"
    }
    """
    Given the "note" exists
    """
    {
       "user_id": "90d78048-f39d-47ab-9d24-3da4d8d1fb23",
       "id": "4c3d0e05-f341-4f03-bae1-f5852bd4487f",
       "title": "My first note",
       "content": "This is my first note"
    }
    """
    When I call "GET" "/v1/users/90d78048-f39d-47ab-9d24-3da4d8d1fb23/notes/4c3d0e05-f341-4f03-bae1-f5852bd4487f"
    Then the status returned should be 200
    And the response should contain the field "id" equal to "4c3d0e05-f341-4f03-bae1-f5852bd4487f"
    And the response should contain the field "title" equal to "My first note"
    And the response should contain the field "content" equal to "This is my first note"

  Scenario: Get note failure note does not exist - response fields validation
    Given the "user" exists
    """
    {
      "id": "90d78048-f39d-47ab-9d24-3da4d8d1fb23",
      "name": "John Doe"
    }
    """
    When I call "GET" "/v1/users/90d78048-f39d-47ab-9d24-3da4d8d1fb23/notes/4c3d0e05-f341-4f03-bae1-f5852bd4487f"
    Then the status returned should be 404

