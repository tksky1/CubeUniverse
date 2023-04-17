import { commonLayout } from "@/utils/commonLayout"
import {
    Button,
    Center,
    Modal,
    NumberInput,
    Stack,
    Switch,
    TextInput,
    Title,
    Tooltip,
    ActionIcon,
    Loader,
    Group,
    Paper,
    Text,
    SimpleGrid,
    ThemeIcon,
    Tabs,
    Divider,
    Image,
    Grid,
    useMantineTheme
} from "@mantine/core";
import { useEffect, useState, useCallback } from "react";
import { checkStorageOpen, createPvc, deletePvc, getKeyList, getPvcInfo, openStorage, updatePvcVolume, searchImg, getObject } from "@/apis";
import { notifications } from "@mantine/notifications"
import { Box } from "@mantine/core";
import { MantineReactTable } from 'mantine-react-table';
import type { MRT_ColumnDef } from 'mantine-react-table';
import { IconX, IconCheck, IconEdit, IconTrash, IconFile, IconFilter, IconPhoto, IconSearch } from "@tabler/icons-react";
import { useDisclosure, useInputState } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { v4 } from "uuid";
import { noti } from "@/utils/noti";
import { useRouter } from "next/router";

interface ObjectPVC {
    name: string,
    namespace: string,
    time: string,
    max_object: string,
    max_size: string
}

export default function ObjectDetail() {
    let { query, replace } = useRouter();
    let [storageOpen, setStorageOpen] = useState<"loading" | boolean>("loading");
    let [pvcs, setPvcs] = useState<ObjectPVC[]>([]);
    let updateStorageOpen = useCallback(() => {
        checkStorageOpen("object")
            .then(e => e.json())
            .then(e => {
                if (e.msg && e.msg[0] === "p") {
                    setStorageOpen(true);
                } else {
                    setStorageOpen(false);
                }
            });
    }, []);
    let updatePvc = useCallback(() => (
        getPvcInfo("object")
            .then(r => r.json())
            .then(r => setPvcs(Object.values(r.data)))
            .catch(r => {
                notifications.show({
                    id: v4(),
                    withCloseButton: true,
                    autoClose: 5000,
                    title: "Error",
                    message: "Something went wrong",
                    color: 'red',
                    icon: <IconX />,
                    loading: false,
                    variant: "outline"
                });
            })), []);
    useEffect(() => {
        updateStorageOpen();
        updatePvc();
    }, []);
    return query.name && query.namespace
        ? <Detail
            name={query.name as string} namespace={query.namespace as string} />
        : (
            <Center h={"100%"}>
                <Button onClick={() => replace("/objectStorage")}>返回对象存储页面</Button>
            </Center>
        )
}

ObjectDetail.getLayout = commonLayout;

interface DetailInterface {
    name: string,
    namespace: string
}

function Detail({ name, namespace }: DetailInterface) {
    let { replace } = useRouter();
    const [activeTab, setActiveTab] = useState<string | null>("files");
    let [keys, setKeys] = useState<string[]>([]);
    let [filterStr, setFilterStr] = useInputState("");
    useEffect(() => {
        getKeyList({ name, namespace })
            .then(e => e.json())
            .then(e => e.code === 200 ? Promise.resolve(e) : Promise.reject(e))
            .then(e => {
                setKeys(e.data.keys);
            })
            .catch(e => {
                setKeys([]);
            })
    }, []);
    return (
        keys.length === 0
            ? (
                <Center h={"100%"}>
                    <Button
                        variant="gradient"
                        gradient={{ from: 'indigo', to: 'cyan' }}
                        onClick={() => replace("/objectStorage")}>返回对象存储页面</Button>
                </Center>
            )
            : (
                <Box
                    p={"lg"}
                    mih={"100%"}
                >
                    <Tabs
                        value={activeTab}
                        onTabChange={val => setActiveTab(val)}
                        mih={"100%"}
                    >
                        <Tabs.List>
                            <Tabs.Tab value="files">Files</Tabs.Tab>
                            <Tabs.Tab value="search">Search</Tabs.Tab>
                        </Tabs.List>
                        <Tabs.Panel value="files" pt="xs">
                            <Paper
                                px={"md"}
                                py={"xs"}
                                shadow="lg"
                                mb={"xs"}>
                                <Group
                                    position="apart">
                                    <TextInput
                                        icon={
                                            <ThemeIcon
                                                sx={{ border: "none" }}
                                                variant="outline" size={"xs"}>
                                                <IconFilter />
                                            </ThemeIcon>}
                                        value={filterStr}
                                        onChange={setFilterStr}
                                    />
                                    <Group>
                                        <Text>命名空间：{namespace}</Text>
                                        <Text>存储卷声明：{name}</Text>
                                    </Group>
                                </Group>
                            </Paper>
                            <Paper shadow="xl" >
                                <SimpleGrid cols={5}>
                                    {keys.filter(k => k.includes(filterStr)).map(k => <FileCard key={k} filename={k} />)}
                                </SimpleGrid>
                            </Paper>
                        </Tabs.Panel>
                        <Tabs.Panel
                            value="search"
                            pt="xs"
                            mih={"100%"}
                        >
                            <Search name={name} namespace={namespace} />
                        </Tabs.Panel>
                    </Tabs>
                </Box>
            )
    )
}

function FileCard({ filename }: { filename: string }) {
    return (
        <Box p={"xs"} >
            <Stack spacing={"xs"}>
                <Center w={"100%"}>
                    <MyIcon str={filename} />
                </Center>
                <Text align="center">{filename}</Text>
            </Stack>
        </Box>
    )
}

function MyIcon({ str }: { str: string }) {
    return (
        (str.endsWith(".jpeg") || str.endsWith(".png") || str.endsWith(".jpg"))
            ? (
                <ThemeIcon
                    size={80}
                    variant="gradient"
                    gradient={{ from: '#ed6ea0', to: '#ec8c69', deg: 35 }}
                >
                    <IconPhoto />
                </ThemeIcon>
            )
            : (
                <ThemeIcon
                    size={80}
                    variant="gradient"
                    gradient={{ from: 'teal', to: 'lime', deg: 35 }}
                >
                    <IconFile />
                </ThemeIcon>
            )
    )
}

interface SearchInterface {
    name: string,
    namespace: string
}

function Search({ name, namespace }: SearchInterface) {
    let [searchStr, setSearchStr] = useInputState("");
    let theme = useMantineTheme();
    let [res, setRes] = useState([]);
    return (
        <Paper
            shadow="xl"
            p={"xs"}
        >
            <Stack>
                <Title
                    order={1}
                    variant="gradient"
                    gradient={{ from: "#0081FA", to: "#CF00A9", deg: 30 }}
                >智能对象查询</Title>
                <Group position="apart">
                    当前对象桶：{name}
                    <Group spacing={"xs"}>
                        <TextInput
                            placeholder="图片关键词"
                            value={searchStr}
                            onChange={setSearchStr}
                        />
                        <ActionIcon
                            onClick={() => searchImg({
                                name,
                                namespace,
                                tag: searchStr
                            }).then(e => e.json())
                                .then(e => e.code === 200
                                    ? Promise.resolve(e)
                                    : Promise.reject(e))
                                .then(e => setRes(e.data.keys))
                                .catch(e => noti("Error", "Error", (e && e.msg)
                                    ? e.msg
                                    : "Search Failed"))}>
                            <IconSearch />
                        </ActionIcon>
                    </Group>
                </Group>
                <Divider />
                {res.length === 0
                    ? (
                        <Center h={"60vh"}>
                            <Title order={5}>
                                输入关键字查询图片
                            </Title>
                        </Center>
                    )
                    : (
                        <SimpleGrid cols={5}>
                            {res.map(r => (
                                <MyPic
                                    name={name}
                                    namespace={namespace}
                                    picKey={r}
                                    key={r}
                                />))}
                        </SimpleGrid>

                    )}
            </Stack>
        </Paper>
    )
}

interface MyPicInterface {
    picKey: string,
    name: string,
    namespace: string
}

function MyPic({ picKey, name, namespace }: MyPicInterface) {
    let [res, setRes] = useState<"loading" | string>("loading");

    useEffect(() => {
        getObject({ name, namespace, key: picKey })
            .then(r => r.json())
            .then(r => setRes(r.data.value0));
    }, []);
    return (
        <Stack>
            {res === "loading"
                ? (
                    <Center h={300}>
                        <Loader />
                    </Center>
                )
                : (
                    <Image
                        height={100}
                        fit="contain"
                        src={`data:image/${picKey.split(".").reduce((p, c) => c)};base64,${res}`}
                    />
                )}
            {/* <p>{`data:image/${picKey.split(".").reduce((p, c) => c)};base64,${res}`}</p> */}
            <Text align="center">{picKey}</Text>
        </Stack>
    )
}

