Feature: Calendar app

  Scenario: Create event
    When I send addEvent request
    Then I want to see event ID in response

  Scenario: Create duplicate
    Given there is event with date "2020-09-15T00:40:08.000Z"
    When I send addEvent request on date "2020-09-15T00:40:08.000Z"
    Then Response has error

  Scenario: Update event
    Given there is event with date "2020-09-28T00:40:08.000Z"
    When I send updateEvent request
    Then Response has NO errors

  Scenario: Update unknown id
    When I send updateEvent request of "b12abe58f-45lk-9673-0924-2s34a2c4f430"
    Then Response has error

  Scenario: Delete event
    Given there is event with date "2020-09-30T00:40:08.000Z"
    When I send deleteEvent request
    Then Response has NO errors

  Scenario: Delete unknown id
    When I send deleteEvent request for "b12abe58f-45lk-9673-0924-2s34a2c4f430"
    Then Response has error

  Scenario: Receive events list
    Given there is event with date "2020-09-20T00:40:08.000Z"
    When I send getList request
    Then I want to see events response