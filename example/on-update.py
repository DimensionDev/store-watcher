#!/usr/bin/env python3
import json
import sys
from datetime import datetime

import requests

API_KEY = "your mailgun api"
DOMAIN = "your mailgun domain"
SENDER = "store-watcher"
TO_ADDRESS = "to email address"
CC_ADDRESS = None
BCC_ADDRESS = None

payload = json.load(sys.stdin)
payload["previous_date"] = payload["previous_date"] or "N/A"

timestamp = datetime.now().strftime("%Y-%m-%d")
subject = "{timestamp} {name} ({platform}) | {previous_version} -> {current_version}"
content = """\
Product: {name}
Platform: {platform}

Previous version: {previous_version}
Previous publish date: {previous_date}

Current version: {current_version}
Current publish date: {current_date}

Link: {link}\
"""

requests.post(
    "https://api.mailgun.net/v3/%s/messages" % DOMAIN,
    auth=("api", API_KEY),
    data={
        "from": "Store Watcher <%s@%s>" % (SENDER, DOMAIN),
        "to": TO_ADDRESS,
        "cc": CC_ADDRESS,
        "bcc": BCC_ADDRESS,
        "subject": subject.format(timestamp=timestamp, **payload),
        "text": content.format_map(payload),
    },
)
