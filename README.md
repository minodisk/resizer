# resizer [ ![Codeship Status for minodisk/resizer](https://img.shields.io/codeship/13aaf0c0-0354-0135-03f6-4ebd009726c2/master.svg?style=flat)](https://app.codeship.com/projects/213308) [![Go Report Card](https://goreportcard.com/badge/github.com/minodisk/resizer)](https://goreportcard.com/report/github.com/minodisk/resizer) [![codecov](https://codecov.io/gh/minodisk/resizer/branch/master/graph/badge.svg)](https://codecov.io/gh/minodisk/resizer) [![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://godoc.org/github.com/minodisk/resizer) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)


## Specification

- Keeps aspect ratio.
- Doesn't scale up, but scale down.
- Reflect orientation tag in EXIF of JPEG to pixels of resized image.
- Drops meta data.

## Installation

Download binary for your environment from [latest release](https://github.com/minodisk/resizer/releases/latest), or `go get` like:

```bash
go get -u github.com/minodisk/resizer/...
```

## Usage

```bash
./resizer
```

### Environment variables

- `ENVIRONMENT`: `development` or `production`. In default `production`
- `GC_PROJECT_ID`: The project ID of Google Cloud.
- `GC_STORAGE_BUCKET`: The bucket name of Google Cloud Storage.
- `GC_SERVICE_ACCOUNT`: The base64-encoded service account JSON of Google Cloud.
- `MYSQL_DATA_SOURCE_NAME`: The data source name of MySQL.
- `ALLOWED_HOSTS`: The allowed hosts of resizing images. Specified with a string joined by `,`. When unspecified, resizer will resize images in any host.
- `MAX_HTTP_CONNECTIONS`

## HTTP(S) API

### Examples

```http:HTTPRequest
curl http://your.host.name/?url=http%3A%2F%2Fexample.com%2Fimage.jpeg&width=800
```

### Endpoint

```http:Endpoint
GET http://your.host.name/
```

### Parameters

- Joint `key=value` parameters with `&`.
  - The `value` should be URL-encoded.

#### `url`

The URL of a resizing image. Required.
The host of the URL should be specified with `hosts` in running option.

#### `width`, `height`

The size of resized image in pixel. In default `0`.

- An integer greater than or equal to `0`.
- Specify `0` to both of `width` and `height` isn't allowed.
- When `width` is `0`. `width` is guessed with `height` and the aspect ratio of the source image .
- When `height` is `0`. `height` is guessed with `width` and the aspect ratio of the source image .
- The specified size is greater than the size of source image, resizer doesn't resize.

#### `method`

How to resize. `normal` or `thumbnail`. Optional. In default `normal`.

- When `width` or `height` is `0`, specified `method` is ignored and resizer resizes `normal`.
  - `width` and `height` shouldn't be specified with `0`, when method is `thumbnail`.
- When specifies `normal`, resizer resizes image to fall into the specified size and doesn't clip.
- When specifies `thumbnail`, resizer resizes image to fill all pixels in the specified size and clips the outer of the specified size.

#### `format`

The format of the resized image. `jpeg` or `png` or `gif`. In default `jpeg`.

#### `quality`

The quality of the resized image as `jpeg`. `0`〜`100`. In default `100`.

- Ignored, when `format` isn't `jpeg`.

### Response

#### Success

- When resizes first time, response resized image data with the code as `2xx`.
- When resizes second (or third or forth) time, response with code as `3xx` and redirects to the storage URL of that the resized image was saved.

#### Error

Response with the code as `4xx`, and the reason will be written in the body.
