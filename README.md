# ftbadge

Generate a dynamic badge displaying a 42 profile summary from a GitHub login.

## Performance & Caching

ftbadge is designed for efficient delivery. While cold starts can introduce a slight delay on the first view due to the hosting environment, the badge is cached for 24 hours.

Our multi-layered caching strategy ensures:

- Fast Load Times: Badges are delivered quickly for an enhanced user experience.
- Reduced External API Calls: Minimizes reliance on external services, improving reliability.
- High Availability: Cached content remains accessible even during potential upstream issues.

This approach ensures your badges are displayed quickly and reliably.

## Usage

```md
<a href="link"><img src="https://ftbadge.cc/login" alt="description"></a>
```

### With size

```md
<a href="link"><img src="https://ftbadge.cc/login" alt="description" width="width" height="width"></a>
```

## Live Exemple

<a href="https://github.com/theoguerin64/ftbadge"><img src="https://ftbadge.cc/tguerin" alt="42 Profile Badge" width="425" height="175"></a>

## Contributing

Contributions are welcome! Whether it’s…

- New SVG badge styles
- Bug reports or feature requests (please open an issue)
- Pull requests with enhancements or fixes

…feel free to fork the repo and submit a PR.
