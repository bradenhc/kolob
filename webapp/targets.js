// Defines all the bundle targets that need to be included in our final output.
//
// The `opts` argument is provided by the common build script and has the following properties:
// - `destdir`: The path to the output destination directory.
// - `release`: True if we are building for release, false otherwise.
//
// The default exported function must return an array of objects that follow the esbuild API config
// object structure. See the following URL for details:
//
//     https://esbuild.github.io/api/#overview
//
export default function(opts) {
  return [
    {
      entryPoints: ["src/index.html"],
      loader: {
        ".html": "copy",
      },
      outdir: opts.destdir,
    },
    {
      entryPoints: ["src/app.js"],
      minify: opts.release,
      bundle: true,
      sourcemap: !opts.release,
      outdir: opts.destdir,
    },
    {
      entryPoints: ["src/styles.css"],
      minify: opts.release,
      outdir: opts.destdir,
    },
  ];
};
