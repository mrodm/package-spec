Feature: Synthetic source
  Support synthetic source

  @3.2.0
  Scenario: Installer leverages synthetic source mode and fields with store true
   Given the "logs_synthetic_mode" package is installed
     And a policy is created with "logs_synthetic_mode" package and "1.0.0-beta1" version
    Then index template "logs-logs_synthetic_mode.synthetic" is configured for "synthetic source mode"
     And index template "logs-logs_synthetic_mode.synthetic" has a field "decision_id" with "store:true"
