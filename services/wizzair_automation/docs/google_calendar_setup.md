# Google Calendar Setup Guide

This guide explains how to set up Google Calendar API credentials and share
your calendar with the Wizz Air Automation.

Basically, you have to:

- Configure the Calendar API in a Google Cloud project to be used with this
  automation.
- Add a service account to the project.
- Download the JSON key for the project.
- Share the Google Calendar with the service account.
- Configure the environment variables for the automation:

  ```sh
  GCAL_SERVICE_ACCOUNT_JSON_PATH="/path/to/service_account.json"
  GCAL_CALENDAR_ID="mycalendar@group.calendar.google.com"
  ```


## Sharing a calendar with a Service Account

The service's account JSON key file has a `client_email` value (e.g.
`wizzair-calendar-sync@<project-id>.iam.gserviceaccount.com`). This is the user
we should share our calendar with.

To do this, open [Google Calendar](https://calendar.google.com/), go to the
configuration, and share your calendar with that email address, giving that
user the permission to "Make changes and see all event details".
