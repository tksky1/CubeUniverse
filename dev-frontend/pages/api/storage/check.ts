import type { NextApiRequest, NextApiResponse } from "next";

type Data = {
    code: number,
    data: null,
    msg: string
}

export default function handler(req: NextApiRequest, res: NextApiResponse<Data>) {
    let action = res.getHeader("X-action");

    if (action === "check") {
        res.status(200).json({
            "code": 200,
            "data": null,
            "msg": "BCS"
        });
    } else {
        res.status(200);
    }
}