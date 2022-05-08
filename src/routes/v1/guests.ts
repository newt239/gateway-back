import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
import { QueryError } from 'mysql2';
const router = express.Router();

router.get('/info/:guest_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userId, res.locals.password);
    const guest_id: string = req.params.guest_id;
    const sql: string = `SELECT * FROM gateway.guest WHERE guest_id='${guest_id}'`;
    connection.query(sql, (err: QueryError, result: any) => {
        if (err) {
            return res.status(400).json(err);
        } else {
            if (result.length === 0) {
                return res.json({
                    message: `${guest_id}というゲストは存在しません。`
                });
            } else {
                return res.json(result[0]);
            };
        };
    });
    connection.end();
});

router.post('/register', verifyToken, (req: express.Request, res: express.Response) => {
    const connection = connectDb(res.locals.userId, res.locals.password);
    // 同じreservationによる複数のguestの登録
    let sql: string = "INSERT INTO gateway.guest (guest_id, guest_type, reservation_id, user_id, register_at, available) VALUES";
    for (const eachRegister of req.body) {
        const guest_id: string = eachRegister.guest_id;
        const reservation_id: string | null = eachRegister.reservation_id;
        const guest_type: string = "student";
        const timestamp = new Date().toLocaleString('ja-JP').slice(0, 19).replace('T', ' ');
        sql += `('${guest_id}', '${guest_type}', '${reservation_id}', '${res.locals.userId}', '${timestamp}', 1),`;
    }
    sql = sql.slice(0, -1) + ";";
    // sessionテーブルへ入場を記録
    sql += "INSERT INTO gateway.session (session_id, exhibit_id, guest_id, enter_at, enter_operation, available) VALUES"
    for (const eachRegister of req.body) {
        const date = new Date();
        const session_id: string = date.getTime().toString(16) + Math.floor(Math.random() * 10).toString(16);
        const guest_id: string = eachRegister.guest_id;
        const timestamp = date.toLocaleString('ja-JP').slice(0, 19).replace('T', ' ');
        sql += `('${session_id}, 'entrance', ${guest_id}', '${timestamp}', '${res.locals.userId}', 1),`;
    }
    sql = sql.slice(0, -1) + ";";
    // 登録したguestの数だけreservation tableのregisteredの数を増加させる
    sql += `UPDATE gateway.reservation SET registered = registered + ${req.body.length} WHERE reservation_id='${req.body[0].reservation_id}';`
    console.log(sql);
    connection.query(sql, (err: QueryError, _result) => {
        if (err) {
            console.log(err);
            if (err.code === "ER_DUP_ENTRY") {
                return res.status(400).json({
                    message: "guest_id is already registered."
                });
            } else {
                return res.status(400).json(err);
            }
        } else {
            return res.json();
        };
    });
    connection.end();
});

router.post('/revoke', verifyToken, (req: express.Request, res: express.Response) => {
    const connection = connectDb(res.locals.userId, res.locals.password);
    const guest_id: string = req.body.guest_id;
    const timestamp = new Date().toLocaleString('ja-JP').slice(0, 19).replace('T', ' ');
    const sql: string = `UPDATE gateway.session SET exit_at='${timestamp}', exit_operation='${res.locals.userId}' WHERE guest_id='${guest_id}' AND exit_at IS NULL;`;
    console.log(sql);
    connection.query(sql, (err: QueryError, _result) => {
        if (err) {
            return res.status(400).json(err);
        } else {
            return res.json();
        };
    });
    connection.end();
});

export default router;