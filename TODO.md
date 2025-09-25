# gcfg's TODO

Want to contribute? Feel free to pick and choose what makes sense for your use case.

## High Impact Features

### Configuration Validation

- [ ] Add a `Validator` interface that extensions can implement
- [ ] Built-in validators for common patterns (required fields, ranges, formats)
- [ ] Schema validation support (maybe JSON Schema or a simple struct-based approach)
- [ ] Validation error aggregation (collect all errors, don't fail on first)

> Note: I'm thinking about using [go-playground/validator](https://github.com/go-playground/validator) for this! And
> maybe auto-validate after `Bind()`?

### YAML Support

- [ ] `YAMLProvider` - because everyone loves YAML configs 🤷
- [ ] Same pattern as JSON provider, just different unmarshaling

> Note: Consider using `gopkg.in/yaml.v3`?

### Performance Optimizations

- [ ] Lazy loading option for large configuration files
- [ ] Memory pool for frequently allocated maps during merging
- [ ] Other optimizations TBA?

> Note: The `Get(key string) any` is *intensionally* not optimized, users should be encouraged to use `Bind()` with
> structured config instead!

## Nice-to-Have Features

### Hot Reloading

- [ ] File watcher integration (maybe `fsnotify`?)
- [ ] Callback system for configuration changes
- [ ] Thread-safe config updates without blocking readers
- [ ] Debouncing for rapid file changes

## Developer Experience Improvements

### Better Error Messages

- [ ] More context in error messages (file paths, line numbers for YAML/JSON)
- [ ] Suggestions for common mistakes ("Did you mean 'database.host'?")
- [ ] Configuration path tracing for debugging

### Observability

- [ ] Optional logging interface for load operations
- [ ] Metrics hooks (how long did loading take? which providers were used?)
- [ ] Debug mode that shows the merge order and final values

### Documentation & Examples

- [ ] More examples in the `/examples` directory
- [ ] Performance benchmarks and best practices guide

## Advanced Features

### Configuration Encryption

- [ ] Encrypted value support (maybe with a simple `gcfg:encrypted:` prefix?)
- [ ] Key management integration (environment variables for keys)
- [ ] Selective field encryption (only sensitive values)

> Note: Is this even needed? And should it be built-in or an extension?

### Remote Configuration

- [ ] HTTP provider for remote configs
- [ ] Consul/etcd/AWS SecretsManager provider for distributed systems
- [ ] S3/cloud storage provider

> Note: These can be implemented as external/separate providers.

### Custom Merge Strategies

- [ ] Array merge strategies (append vs replace vs merge by key)
- [ ] Custom merger interface for complex types
- [ ] Merge conflict detection and resolution

### Type Safety Improvements

- [ ] Make the `Config` struct generic or add a generic variant!
- [ ] Compile-time configuration struct validation
- [ ] Better error messages for type mismatches

## Polish & Maintenance

### Testing & Quality

- [ ] Benchmark tests for performance tracking
- [ ] Fuzzing tests for robustness
- [ ] Integration tests with real applications
- [ ] Race conditions detection tests
- [ ] Memory leak detection tests

### Community

- [ ] Contribution guidelines
- [ ] Issue templates
- [ ] Automated security scanning
- [ ] Dependency update automation
