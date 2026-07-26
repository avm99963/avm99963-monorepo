# Wizz Air automation

An automation that, once I receive an itinerary confirmation email from Wizz
Air, does the following:

- Adds events to my calendar for those flights.
- Adds check-in reminders to my to do list.

## Spin up this project

First of all, create an email account (in Gmail, or wherever) that will receive
emails by Wizz Air. This account should be used exclusively with this
automation. Take note of the credentials and info that can be used to query
email via IMAP, you will set these up in the environment variables.

Then, create a Google Calendar where you will be saving the flight events. In
Google Cloud, create a project and [set up the Google Calendar
API](docs/google_calendar_setup.md).

You should also have a task list for this purpose set up in software like
Nextcloud, and available via CalDAV.

Now you're ready to make a copy of [.env.sample](.env.sample) and set your
environment variables to the appropriate values to wire everything together.

You can then run the OCI image available at
`ghcr.io/avm99963/wizzair-automation` with these environment variables in
Docker, Kubernetes, or wherever you want to! Just spin up an instance with that
image and keep it running. At startup and every 12 hours, it will perform the
automation.

## Configuration

| Environment variable | Description | Default |
| :--- | :--- | :--- |
| `POLL_INTERVAL` | Polling frequency duration | `12h` |
| `TIMEZONE` | Fallback timezone string | `Europe/Warsaw` |
| `DRY_RUN` | Dry-run mode (`true`/`false`).<br>When `true`, emails will be downloaded but the automation will not actually perform any change. Instead, logs will be written with the actions that would be taken. | `false` |
| `IMAP_SERVER` | IMAP server hostname | *Required* |
| `IMAP_PORT` | IMAP TLS/STARTTLS port | `993` |
| `IMAP_USERNAME` | IMAP account username/email | *Required* |
| `IMAP_PASSWORD` | IMAP account password | *Required* |
| `LABEL_PROCESSED` | IMAP label/folder for processed itineraries | `Processed` |
| `LABEL_IGNORED` | IMAP label/folder for non-itinerary emails | `Processed (ignored)` |
| `GCAL_CALENDAR_ID` | Target Google Calendar ID | `primary` |
| `GCAL_SERVICE_ACCOUNT_JSON_PATH` | Absolute path to GCP Service Account JSON key | *Required* |
| `CALDAV_BASE_URL` | CalDAV task calendar URL | *Required* |
| `CALDAV_USERNAME` | Username | *Required* |
| `CALDAV_PASSWORD` | Password | *Required* |
