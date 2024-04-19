//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//

// This Node.js script can bundle web applications using esbuild. It supports watching and serving
// the generated files during developent. It can also proxy an API server to handle API
// requests sent to the esbuild development server.
//
// Run the script with '-h' for additional information

import fs from "node:fs";
import http from "node:http";

import esbuild from "esbuild";
import { program } from "commander";

// Parse command-line arguments
//
const DESCRIPTION = [
  "Bundles JavaScript source code targeting either the browser, Node.js, or both. This script ",
  "expects there to be a `targets.js` file in the current working directory that defines all of ",
  "the bundle targets using the esbuild build API config object structure.\n\n",
  "",
  "For more information about the build API config object structure, see the corresponding ",
  "esbuild documetation: https://esbuild.github.io/api/#overview.\n\n",
  "",
  "This script also uses esbuild's native development server to serve bundled source during ",
  "development. It will even watch the source files for changes and automatically rebuild when ",
  "they are detected.\n\n",
  "",
  "In addition to this, the script also supports forwarding requets that cannot be handled by the ",
  "development server to another server using the '-f,--forward' option. This is especially ",
  "useful when you are building a web application that uses some service API and you want to ",
  "sidestep the CORS errors you'll likely get trying to hit REST API endpoints on the service ",
  "that are not coming from the same host/port the development server is running on.",
].join("");

program.description(DESCRIPTION);
program.option(
  "-f,--forward <value>",
  "Set the address to forward requests that can't be handled by the developent server to.",
);
program.option("-o,--output <dir>", "Set the output directory (default 'dist')", "dist");
program.option("-p,--port <value>", "use the specified port for the dev server", 1234);
program.option("-r,--release", "Optimie the output file size.");
program.option("-s,--serve", "Build and start a dev server.");
program.option("-w,--watch", "Build and watch files for changes.");

program.configureHelp({ helpWidth: 88 });
program.parse();

const opts = program.opts();

// Define the async function that will run the development server and also forward unhandled
// requests to another service if the script is configured to use one. The ability to forward
// requests on allows us to run an API server on the same host and proxy the requests to it so
// we don't get CORS errors in the browser when developing apps.
//
async function serve(ctx, servedir, devPort, forwardAddress) {
  const esAddress = await ctx.serve({ servedir });
  console.info(`Started esbuild server on ${esAddress.host}:${esAddress.port}`);

  const proxy = http.createServer((req, res) => {
    const forwardRequest = (host, port, path) => {
      const options = {
        host,
        port,
        path,
        method: req.method,
        headers: req.headers,
      };

      const proxyReq = http.request(options, (proxyRes) => {
        if (proxyRes.statusCode === 404 && forwardAddress !== undefined) {
          // if esbuild 404s the request, assume that our REST API server is going to
          // hadle the request because it is likely an API call. Forward it.
          const [fhost, fport] = forwardAddress.split(":");
          return forwardRequest(fhost, parseInt(fport), req.url);
        }

        // Otherwise esbuild handled it like a chap, so proxy the response back.
        res.writeHead(proxyRes.statusCode, proxyRes.headers);
        proxyRes.pipe(res, { end: true });
      });

      req.pipe(proxyReq, { end: true }).on("error", () => {
        console.error(`Failed to proxy request: ${req.method} ${req.url}`);
        res.statusCode = 404;
        res.end();
      });
    };

    // When we're called pass the request right through to esbuild
    forwardRequest(esAddress.host, esAddress.port, req.url);
  });

  proxy.listen(devPort);
  console.info(`Started dev server on 0.0.0.0:${devPort}`);
  if (forwardAddress) {
    console.info(`Requests unservable by esbuild will be forwarded to ${forwardAddress}`);
  }
}

// Define our main function, which we immediately call. We wrap all our build instructions inside
// of an async IIFE so that we can use the async/await capabilities of JavaScript for improved
// code readability.
//
(async function () {
  const tasks = [];
  const contexts = [];

  try {
    // Make sure the output directory exists.
    await fs.promises.mkdir(opts.output).catch(() => {});

    // Import the target definitions. Defining it before using it enables us to use the same
    // configuration for both develompent and production builds.
    const targets = await import(process.cwd() + "/targets.js");

    // Create tasks to build all our targets. We build and watch for file changes if the user
    // provides the flag that enables this behavior.
    for (const t of targets.default({ destdir: opts.output, release: opts.release })) {
      if (opts.watch) {
        const ctx = await esbuild.context(t);
        contexts.push(ctx);
        tasks.push(ctx.watch());
      } else {
        tasks.push(esbuild.build(t));
      }
    }

    // Start the development server if the user used the flag that enables this behavior
    if (opts.serve) {
      const ctx = await esbuild.context({});
      contexts.push(ctx);
      await serve(ctx, opts.output, parseInt(opts.port), opts.forward);
    }

    // We started a bunch of async tasks to build/watch/serve our project. Wait for all the
    // tasks to complete.
    await Promise.all(tasks);
  } catch (e) {
    // One of the async tasks failed. Log the error and stop all other tasks.
    console.error(e);
    for (const c of contexts) {
      await c.dispose();
    }
    process.exitCode = 1;
  }
})();
