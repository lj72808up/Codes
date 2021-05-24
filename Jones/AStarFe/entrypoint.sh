#!/bin/bash
sed -i "s@BASE_API_REPLACE@$BASE_API_URL@g" /usr/share/nginx/html/static/js/*
