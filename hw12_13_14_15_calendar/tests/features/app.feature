Feature: Calendar app

  Scenario: Create event
    When I send addEvent request
    Then I want to see event ID in response

  Scenario: Create duplicate
    Given I have event on date "2020-09-15T10:30:00.000Z"
    When I send addEvent request on date "2020-09-15T10:30:00.000Z"
    Then Response has error

  Scenario: Create duplicate
     When I send addEvent request on today
     Then I want to see event ID in response
     Then Response has error
     And I want receive event notification

  Scenario: Update event
    Given there is event with date "2020-09-15T10:30:00.000Z"
    When I send updateEvent request
    Then Response has NO errors

  Scenario: Update unknown id
    When I send updateEvent request for "b12abe58f-45lk-9673-0924-2s34a2c4f430"
    Then I get error response

  Scenario: Delete event
    Given there is event with date "2020-09-15T10:30:00.000Z"
    When I send deleteEvent request
    Then Response has NO errors

  Scenario: Delete unknown id
    When I send deleteEvent request for "b12abe58f-45lk-9673-0924-2s34a2c4f430"
    Then Response has NO errors

  Scenario: Receive events list
    When I send getList request
    Then I want to see events response