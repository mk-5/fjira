Feature: project selection
  The center point of Jira is Jira Project. Fjira user needs to select an project
  in order to do some actions with tickets.

  Scenario: Open fjira and select project
    Given projects fuzzy find is up&running
    When project is selected
    Then fjira should open project view

  Scenario: Open project directly from terminal
    Given CLI argument with project key is present
    When fjira started
    Then fjira should open project view
