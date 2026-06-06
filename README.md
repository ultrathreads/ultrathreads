# UltraThreads

![Screenshot](screenshot.png)

`UltraThreads` is an open-source, lightweight web forum that renders posts in threaded view.

## Features

- data storage in a MySQL or SQLite database
- nodes
- three possible views of posts

## Project Structure

### api

> Developed with `Go`, providing RESTful-style APIs.

*Tech Stack*
- gin (https://github.com/gin-gonic/gin) Go web framework
- JWT (https://github.com/appleboy/gin-jwt) JWT Middleware for Gin framework
- gorm (http://gorm.io/) ORM framework for Go language

### web

> Frontend page rendering service implemented based on `Next.js`.

*Tech Stack*
- Next.js (https://nextjs.org) The React Framework for the Web

## Installation

The project is still in early development. Refer to the README files of the `api` and `web` modules for installation and development instructions.

## License
UltraThreads is open-sourced software licensed under the [GNU General Public License version 3](https://opensource.org/license/gpl-3.0)

