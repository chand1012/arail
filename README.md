# Arail

Arail is an Autonomous Research AI and Intelligent Learning system. To install clone the repo and run `just build` .

## Building

```sh
git clone https://github.com/chand1012/arail.git
cd arail
just build
```

The executable can be found at `bin/arail` .

## Usage

Set your configuration with `arail config` .

```sh
arail config -k $OPENAI_API_KEY -m gpt-4 # or gpt-3.5-turbo
```

After the config is set, we can start using the CLI.

```sh
arial research "What is the meaning of life?"
```

Commands and their options can be found with `arail help` .
