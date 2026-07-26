# Nextcloud CalDAV Tasks Setup Guide

In case you're using Nextcloud, like I am, this guide explains how to configure
the Nextcloud CalDAV integration.

## Obtain CalDAV Task Calendar URL

Suppose that Nextcloud's base DAV endpoint is
`https://nextcloud.example.com/remote.php/dav`.

To get the specific URL for your Tasks calendar:

1. Go to your tasks list in the Nextcloud website.
1. Click the options icon next to the list and click "Copy link".
   - Full URL structure: `https://nextcloud.example.com/remote.php/dav/calendars/<USERNAME>/<CALENDAR_SLUG>/`
   - Example: `https://nextcloud.example.com/remote.php/dav/calendars/avm99963/personal/`

This is what you would enter in `CALDAV_BASE_URL`.
