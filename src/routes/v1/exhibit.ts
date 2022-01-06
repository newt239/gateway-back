import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
const router = express.Router();

router.get('/heat', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.params.guest_id;
    const sql: string = `SELECT * FROM gateway.guest WHERE guest_id='${guest_id}'`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: {}
            });
        };
    });
});

router.get('/info/', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const sql: string = `SELECT exhibit_id, exhibit_name FROM gateway.exhibit`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.status(400).json(err);
        } else {
            return res.json({
                status: "success",
                data: result
            });
        };
    });
});

router.get('/info/:exhibit_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const exhibit_id: string = req.params.exhibit_id;
    const sql: string = `SELECT * FROM gateway.exhibit WHERE exhibit_id='${exhibit_id}'`;
    connection.query(sql, function (err: any, result: any) {
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
});

router.get('/crowd/', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const sql: string = `SELECT timestamp(DATE_FORMAT(timestamp, '%Y-%m-%d %H:00:00')) AS timestamp, COUNT(*) AS count FROM gateway.activity GROUP BY DATE_FORMAT(timestamp, '%Y%m%d%H');`;
    connection.query(sql, function (err: any, result: any) {
        console.log(result);
        const ctime = result[0].timestamp;
        // TODO: 1件も記録がない時間も配列に含めて返す
        const data = result.reduce((accumulator: any, element: any) => {
            const pushList = [];
            const timestamp = element.timestamp;
            while (true) {
                console.log(timestamp, ctime);
                if (timestamp === ctime) {
                    pushList.push({ time: timestamp, count: element.count });
                    break;
                } else {
                    const addtime = ctime.setHours(ctime.getHours() + 1);
                    pushList.push({ time: addtime, count: 0 });
                    break;
                };
            };
            console.log(accumulator);
            return accumulator.concat([pushList]);
        }, []);
        console.log(data);
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: data
            });
        };
    });
});

router.get('/current/', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const sql: string = `SELECT exhibit_id, count(*) FROM gateway.guest GROUP BY exhibit_id WHERE exhibit_id IS NOT NULL;`;
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
});

router.get('/current/:exhibit_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const exhibit_id: string = req.params.exhibit_id;
    const sql: string = `SELECT guest_id AS id, guest_type FROM gateway.guest WHERE exhibit_id='${exhibit_id}';`;
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
});

module.exports = router;