import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
import { QueryError } from 'mysql2';
import triageError from '@/error';
const router = express.Router();

router.get('/heat', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.params.guest_id;
    const sql: string = `SELECT * FROM gateway.guest WHERE guest_id='${guest_id}'`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            return triageError(err);
        } else {
            return res.json({
                status: "success",
                data: {}
            });
        };
    });
    connection.end();
});

router.get('/info/', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const sql: string = `SELECT exhibit_id, exhibit_name FROM gateway.exhibit`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            return res.status(400).json(err);
        } else {
            return res.json({
                status: "success",
                data: result
            });
        };
    });
    connection.end();
});

router.get('/info/:exhibit_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const exhibit_id: string = req.params.exhibit_id;
    const sql: string = `SELECT * FROM gateway.exhibit WHERE exhibit_id='${exhibit_id}'`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: {
                    exhibit_id: result[0].exhibit_id,
                }
            });
        };
    });
    connection.end();
});

router.get('/enter-chart/:exhibit_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const exhibit_id: string = req.params.exhibit_id;
    const day: string = req.query.day as string;
    const sql: string = `SELECT timestamp(DATE_FORMAT(enter_at, '%Y-%m-%d %H:00:00')) AS time, COUNT(*) AS count FROM gateway.session WHERE exhibit_id='${exhibit_id}' AND DATE(enter_at) = '${day}' GROUP BY DATE_FORMAT(enter_at, '%Y%m%d%H');`;
    console.log(sql);
    connection.query(sql, function (err: any, result: any) {
        console.log(result);
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: result
            });
        };
    });
    connection.end();
});

router.get('/current/', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const sql: string = `SELECT session.exhibit_id, COUNT(*) FROM session INNER JOIN guest ON session.guest_id = guest.guest_id GROUP BY session.exhibit_id;`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: result
            });
        };
    });
    connection.end();
});

router.get('/current/:exhibit_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const exhibit_id: string = req.params.exhibit_id;
    const sql: string = `SELECT session.guest_id AS id, guest_type, enter_at FROM session INNER JOIN guest ON session.guest_id = guest.guest_id WHERE session.exhibit_id='${exhibit_id}' AND session.exit_at IS NULL;`;
    connection.query(sql, function (err: any, result: object[]) {
        if (err) {
            res.json(err);
        } else {
            res.json({
                status: "success",
                data: result
            });
        };
    });
    connection.end();
});

module.exports = router;