import { createApp } from "./app.js";
import { config } from "./config/index.js";
import { migrationService } from "./infrastructure/database/migrations.js";
import { logger } from "./infrastructure/logging/logger.js";
import { db } from "./infrastructure/database/connection.js";

const app = createApp();

async function startServer() {
  try {
    logger.info("Initializing database...");
    await migrationService.runMigrations();

    const server = app.listen(config.server.port, () => {
      logger.info(
        {
          port: config.server.port,
          env: config.server.nodeEnv,
        },
        "Auth Service running"
      );
    });

    let isShuttingDown = false;

    const gracefulShutdown = (signal: string) => {
      if (isShuttingDown) {
        return;
      }
      isShuttingDown = true;

      logger.info({ signal }, "Shutting down gracefully");

      const forceExitTimer = setTimeout(() => {
        logger.error(
          "Could not close connections in time, forcefully shutting down"
        );
        process.exit(1);
      }, 10000);

      server.close((err) => {
        if (err) {
          logger.error({ err }, "Error closing HTTP server");
        } else {
          logger.info("HTTP server closed");
        }

        void db
          .close()
          .catch((dbErr: unknown) => {
            logger.error({ err: dbErr }, "Error closing database pool");
          })
          .finally(() => {
            clearTimeout(forceExitTimer);
            process.exit(err ? 1 : 0);
          });
      });
    };

    process.on("SIGTERM", () => {
      gracefulShutdown("SIGTERM");
    });
    process.on("SIGINT", () => {
      gracefulShutdown("SIGINT");
    });

    process.on("unhandledRejection", (reason) => {
      logger.error({ reason }, "Unhandled Promise Rejection");
    });

    process.on("uncaughtException", (error) => {
      logger.error({ error }, "Uncaught Exception");
      gracefulShutdown("UNCAUGHT_EXCEPTION");
    });
  } catch (error) {
    logger.error({ error }, "Failed to start server");
    process.exit(1);
  }
}

void startServer();

