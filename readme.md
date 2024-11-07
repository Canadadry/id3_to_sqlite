# ID3 Tag Management Tool with SQLite Integration

This project provides a tool to extract ID3 tags from media files (such as `.mp3`, `.mp4`, etc.) and store the extracted metadata in a SQLite database. After processing and potentially manipulating the database, the tool also allows updating the original files with any changes made.

## Prerequisites

Ensure you have Go installed on your machine before running the project.

## Installation

To build and run this project:

1. Clone the repository:

   ```bash
   git clone <repository-url>
   ```

2. Navigate to the project directory:

   ```bash
   cd <project-directory>
   ```

3. Install dependencies (if any):

   ```bash
   go mod tidy
   ```

## Project Overview

The tool works in two phases:

1. **Dumping ID3 Tags into a SQLite Database**: This phase scans specified directories for supported media files, extracts ID3 tags (or similar metadata), and saves the extracted data into a SQLite database for further manipulation.

2. **Updating Files from the Database**: After making changes to the metadata in the database, this phase allows updating the original files based on the modified data.

## Commands

### Command 1: Dump Metadata from Files into a SQLite Database

This command scans the specified input directory for media files, extracts metadata (such as ID3 tags), and stores it in a SQLite database.

**Usage:**

```bash
go run main.go dump -i <input-directory> -o <output-database> -c "<columns>" -e <file-extension>
```

**Options:**
- `-i` : Input directory containing media files to be scanned.
- `-o` : Output SQLite database file.
- `-c` : Metadata columns to extract, specified as a semicolon-separated string (e.g., "Album/Movie/Show title;Artist;Title;Year").
- `-e` : File extension filter (e.g., `.mp3`, `.mp4`).

**Example:**

```bash
go run main.go dump -i testdata/ -o db.sqlite -c "Album/Movie/Show title;Artist;Title;Year" -e ".mp3"
```

This example scans the `testdata/` directory for `.mp3` files, extracts specified metadata tags, and stores them in the `db.sqlite` database.

### Command 2: Save and Update Files from the SQLite Database

This command reads data from the specified SQLite database and updates the original media files based on any changes made to their metadata within the database.

**Usage:**

```bash
go run main.go save -i <input-database>
```

**Options:**
- `-i` : Input SQLite database file containing modified metadata.

**Example:**

```bash
go run main.go save -i db.sqlite
```

This example reads data from `db.sqlite` and updates the original media files with any changes made to their metadata in the database.

## Column available

See [here](https://github.com/n10v/id3v2/blob/v2.1.4/v2/common_ids.go) for a full list of available column that can be extracted or set in file
