import { ReactElement } from "react";
import { Layout, Menu } from "antd";
import styles from "@/styles/components/MainLayout.module.scss"
import { useRouter } from "next/router";
const {Header, Sider, Content} = Layout;

type SingleLayerEntry = {
    title: string,
    path: string
}

type DoubleLayersEntry = {
    title: string,
    childrenEntries: SingleLayerEntry[]
}

export type Entry = SingleLayerEntry | DoubleLayersEntry;

function isSingleLayerEntry(entry: Entry): entry is SingleLayerEntry {
    return (entry as SingleLayerEntry).path === undefined;
}

interface MainLayout {
    children: ReactElement
}

let items = [
    {
        title: "page1",
        key: "/page1",
        label: "page1"
    },
    {
        title: "page2",
        key: "/page2",
        label: "page2"
    },
    {
        children: [
            {
                title: "page3",
                key: "/nest/page3",
                label: "page3"
            },
            {
                title: "page4",
                key: "/nest/page4",
                label: "page4"
            },
        ],
        key: "/nest",
        label: "nest"
    },
    {
        type: "group" as const,
        label: "group",
        children: [
            {
                title: "page5",
                key: "/group/page5",
                label: "page5"
            },
            {
                title: "page6",
                key: "/group/page6",
                label: "page6"
            },
        ]
    }
]

export default function MainLayout({children}: MainLayout) {
    let router = useRouter();
    return (
        <div className={styles.outerContainer}>
            <div className={styles.header}>
                Header
            </div>
            <div className={styles.bottomContainer}>
                <div className={styles.sider}>
                    <Menu 
                        theme="light"
                        selectedKeys={[router.asPath]}
                        items={items}
                        mode="inline"
                        onClick={(clickedItem) => router.push(clickedItem.key)}
                        style={{
                            height: "100%",
                            width: "100%",
                        }}
                    />
                </div>
                <div className={styles.contentContainer}>
                    {children}
                </div>
            </div>
        </div>
    )
}