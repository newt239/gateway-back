import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
import { QueryError } from 'mysql2';
const router = express.Router();

router.post('/create', verifyToken, function (req: express.Request, res: express.Response) {
  const connection = connectDb(res.locals.userId, res.locals.password);
  const userId: string = req.body.userId;
  let sql = `CREATE USER '${userId}'@'localhost' IDENTIFIED BY '${req.body.password}'; `;
  switch (req.body.userType) {
    case "moderator":
      sql += `GRANT ALL ON *.* TO '${userId}'@'localhost' WITH GRANT OPTION; `;
      break;
    case "executive":
      sql += `GRANT INSERT, UPDATE, SELECT ON gateway.* TO '${userId}'@'localhost'; `;
      break;
    case "exhibit":
      sql += `GRANT INSERT, UPDATE, SELECT ON gateway.* TO '${userId}'@'localhost'; `;
      break;
    case "analysis":
      sql += `GRANT SELECT ON gateway.* TO '${userId}'@'localhost'; `;
      break;
    default:
      return ({ status: "error", message: "userType wes incorrect." });
  };
  sql += `FLUSH PRIVILEGES;`;
  sql += `INSERT INTO gateway.user (user_id, display_name, user_type, created_by, available) VALUES ('${userId}', '${req.body.displayName}', '${req.body.userType}', '${res.locals.userId}', 1);`
  connection.query(sql, function (err: QueryError, result: any) {
    if (err) {
      if (err.code === "ER_CANNOT_USER") {
        return res.json({
          status: "error",
          message: "ユーザーを作成できませんでした。同じidのユーザーが存在する可能性があります。"
        });
      } else {
        console.log(err);
        return res.json({
          status: "error",
          message: err.message
        })
      }
    } else {
      return res.json({
        status: "success",
        data: result
      });
    }
  });
  connection.end();
});

router.get('/created-by-me', verifyToken, function (req: express.Request, res: express.Response) {
  const connection = connectDb(res.locals.userId, res.locals.password);
  const sql: string = `SELECT user_id, display_name, user_type FROM gateway.user WHERE created_by='${res.locals.userId}'`;
  connection.query(sql, function (err: QueryError, result: any) {
    if (err) {
      return res.json({
        status: "error",
        message: err.message
      });
    } else {
      return res.json({
        status: "success",
        data: result
      });
    }
  });
  connection.end();
});

router.post('/delete-user', verifyToken, function (req: express.Request, res: express.Response) {
  const connection = connectDb(res.locals.userId, res.locals.password);
  const sql: string = `DROP USER '${req.body.userId}'@'localhost'; DELETE FROM gateway.user WHERE user_id='${req.body.userId}' AND created_by='${res.locals.userId}'; `;
  connection.query(sql, function (err: QueryError, result: any) {
    if (err) {
      return res.json({
        status: "error",
        message: err.message
      });
    } else {
      return res.json({
        status: "success",
        userId: req.body.userId
      });
    }
  });
  connection.end();
});

module.exports = router;