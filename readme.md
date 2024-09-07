# Duplicate Image Finder

The Duplicate Image Finder is a Go program that identifies duplicate images in a folder by calculating a hash of the image contents. By comparing the hash of each image, it effectively detects identical files, even if they have different filenames.

## How It Works

The program processes images in the specified folder, computes a hash for each image file, and compares the hashes to detect duplicates. It supports common image formats like JPEG, PNG, and more, and efficiently identifies files with identical contents.

# Features

- Detects duplicate images based on content, not file names.
- Fast and efficient hashing of images.
- Supports various image formats (JPEG, PNG).
- Outputs a list of duplicates in a user-friendly format.

## Prerequisites

To run this program, you need:

- Go installed on your machine (version 1.22.5 or higher).
- A folder containing images to scan for duplicates.

## Installation

- Clone the repository:

```
git clone https://github.com/abuzaforfagun/duplicate-image-finder.git
cd duplicate-image-finder
```

- Install dependencies: `go mod tidy`

## How to Run

_Option 1:_
To run the program from your terminal: `go run . "{YOUR FOLDER PATH}"`
Replace {YOUR FOLDER PATH} with the path to the folder you want to scan for duplicate images.
For example: `go run . "C:\Personal Interest\Golang\duplicate-image-finder\images"`

_Option 2:_
Install the program from your terminal: `go install`
Open the terminal from the folder that need to be checked, and run `duplicate-image-finder`.

## Contribution

If you want to contribute:

- Fork the repository.
- Create a feature branch.
- Commit your changes.
- Open a pull request.
- Feel free to suggest improvements, report issues, or request additional features.
