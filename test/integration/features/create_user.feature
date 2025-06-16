@all @create_user
Feature: Create user
  In order to create a user
  The name should be passed as a request body parameter with minimum 3 characters
  The user must be created in the database
  I need to receive the created user back with the id and the name
  I need to be able to see the payer persisted on the DB with the correct name

  Background:
    Given the header is empty

  Scenario: Create user error when name is less than 3 characters - response fields validation
    Given the header is empty
    When I call "POST" "/v1/users" with body
    """
    {
      "name": "Ab"
    }
    """
    Then the status returned should be 400
    And the response should contain the field "status_code" equal to "400"
    And the response should contain the field "messages.0" equal to "'name' should be greater in length"
    And the response should contain the field "type" equal to "Validation Error"

  Scenario: Create user success - response fields validation
    Given the header is empty
    When I call "POST" "/v1/users" with body
    """
    {
      "name": "John Doe"
    }
    """
    Then the status returned should be 201
    And the response should contain the field "name" equal to "John Doe"
    And the response should contain the field "id" equal to "not nil"

  Scenario: Create user success - db fields validation
    Given the header is empty
    And the db should contain 0 objects in the "users" table
    When I call "POST" "/v1/users" with body
    """
    {
      "name": "John Doe"
    }
    """
    Then the status returned should be 201
    And the db should contain 1 objects in the "users" table
    And the db should contain the "user" with the "name" column value "John Doe" colum "id" equal to "not nil"