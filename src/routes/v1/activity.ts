import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
const router = express.Router();

router.post('/:activity_type', verifyToken, function (req: express.Request, res: express.Response) {
    const activity_type: string = req.params.activity_type;
    if (["enter", "exit", "pass"].indexOf(activity_type, -1)) {
        const connection = connectDb(res.locals.userid, res.locals.password);
        const guest_id: string = req.body.guest_id;
        const exhibit_id: string = req.body.exhibit_id;
        const date = new Date();
        const activity_id: string = date.getTime().toString(16) + Math.floor(Math.random() * 10).toString(16);
        const timestamp = date.toLocaleString('ja-JP').slice(0, 19).replace('T', ' ');
        let sql: string = `UPDATE gateway.guest SET exhibit_id=${activity_type === "enter" ? "'" + exhibit_id + "'" : "NULL"} WHERE guest_id='${guest_id}';`
        if (activity_type === "enter") {
            sql += `INSERT INTO gateway.session (session_id, exhibit_id, guest_id, enter_at, enter_operation, available) VALUES ('s${activity_id}', '${exhibit_id}', '${guest_id}', '${timestamp}', '${res.locals.userid}', 1);`
        } else if (activity_type === "exit") {
            sql += `UPDATE gateway.session SET exit_at='${timestamp}', exit_operation='${res.locals.userid}' WHERE guest_id='${guest_id}' AND exhibit_id='${exhibit_id}' AND exit_at IS NULL;`;
        } else if (activity_type === "pass") {
            sql += `INSERT INTO gateway.session (session_id, exhibit_id, guest_id, enter_at, enter_operation, exit_at, exit_operation, available) VALUES ('s${activity_id}', '${exhibit_id}', '${guest_id}', '${timestamp}', '${res.locals.userid}', '${timestamp}', '${res.locals.userid}', 1);`
        }
        connection.query(sql, function (err: any, result: any) {
            if (err) {
                console.log(err);
                return res.json(err);
            } else {
                return res.json({
                    status: "success",
                    data: { activity_id: activity_id }
                });
            };
        });
        connection.end();
    } else {
        res.json({ status: "error", message: "you posted invaild activity_type" });
    };
});

module.exports = router;