import { ReactNode, useState, useEffect } from "react"
import { DataContext } from "./DataContext"
import useWebSocket, { ReadyState } from "react-use-websocket";

export default function DataProvider({ children }: { children: ReactNode }) {
    let [myData, setMyData] = useState(null);
    let { lastMessage } = useWebSocket("ws://control-backend.cubeuniverse.svc.cluster.local:30401/api/storage/pvcws");
    useEffect(() => {
        if (lastMessage) {
            let newData = JSON.parse(lastMessage.data) as typeof myData;
            console.log(JSON.parse(lastMessage.data));
            setMyData(newData);
        }
    }, [lastMessage]);
    return (
        <DataContext.Provider value={myData}>
            {children}
        </DataContext.Provider>
    )
}