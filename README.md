# tinycounter ![example image hit counter](https://count.mat.dog/rainy/mn6)

`tinycounter` is a web service that provides customizable hit counters as images. It uses Redis for storage and supports different styles for the counter images.

## Features

- No client-side JavaScript required for end users.
  - Embed counters with an HTML `<img>` tag; unique visits are tracked by IP address with a cooldown for incrementing counters.
- Retrieve hit counter images for different users.
- Store and retrieve counter images in various styles.
- Uses Redis for fast and efficient storage.

## Running with Docker

You can run `tinycounter` using Docker and Docker Compose. Here's how to get started in a development environment:

1. Clone the repository:

   ```bash
   git clone https://github.com/mn6/tinycounter.git
   cd tinycounter
   ```

2. Create a `.env` file based on the `.env.example` file:

   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your desired configuration. Update the `users.yaml` if you want to set up whitelists or blacklists for user creation.

4. Start the services using Docker Compose:

   ```bash
   docker compose up -d
   ```

5. Access your counter image at `http://localhost:{APP_PORT}/{style}/{counter id/name}`.

## Extending Styles

To add a new style for the counter images, follow these steps:

1. Create a new directory under `resources/styles/` with the name of your style (e.g., `my-style`).
2. Add digit images (`0.png` to `9.png`) for your style in the new directory. Optionally, add `pre.png` and `suf.png` for prefix and suffix images. All images should be in PNG format and have the same height. Images will be stitched together horizontally to form the final counter image.
3. Update the `styles.yaml` configuration file to include your new style. Width and height should match the dimensions of your digit images. SuffixWidth and PrefixWidth should match the widths of your suffix and prefix images, respectively.
4. Restart the `tinycounter` service to load the new style.
