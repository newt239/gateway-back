import express from 'express';
var app = express();
var jwt = require('jsonwebtoken');

function verifyToken(req: express.Request, res: express.Response, next: express.NextFunction) {
    const authHeader = req.headers["authorization"];
    //HeaderにAuthorizationが定義されているか
    if (authHeader !== undefined) {
        //Bearerが正しく定義されているか
        if (authHeader.split(" ")[0] === "Bearer") {
            try {
                const token = jwt.verify(authHeader.split(" ")[1], 'my_secret');
                //tokenの内容に問題はないか？
                //ここでは、usernameのマッチと有効期限をチェックしているが必要に応じて発行元、その他の確認を追加
                //有効期限はverify()がやってくれるみたいだがいちおう・・・
                if (token.username === "hoge" && Date.now() < token.exp * 1000) {
                    console.log(token);
                    //問題がないので次へ
                    next();
                } else {
                    res.json({ error: "auth error" })
                }
            } catch (e: any) {
                //tokenエラー
                console.log(e.message);
                res.json({ error: e.message })
            }
        } else {
            res.json({ error: "header format error" });
        }
    } else {
        res.json({ error: "header error" });
    }
}
export default verifyToken;