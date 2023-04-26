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
    useMantineTheme,
    Highlight
} from "@mantine/core";
import { useEffect, useState, useCallback } from "react";
import { checkStorageOpen, createPvc, deletePvc, getPvcInfo, openStorage, updatePvcVolume } from "@/apis";
import Image from "next/image";
import { notifications } from "@mantine/notifications"
import { Box } from "@mantine/core";
import { MantineReactTable } from 'mantine-react-table';
import type { MRT_ColumnDef } from 'mantine-react-table';
import { IconX, IconCheck, IconEdit, IconTrash, IconListDetails } from "@tabler/icons-react";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { v4 } from "uuid";
import { noti } from "@/utils/noti";
import { useRouter } from "next/router";
import objectStorage from "../public/obj.webp"

interface ObjectPVC {
    name: string,
    namespace: string,
    time: string,
    max_object: string,
    max_size: string
}

export default function ObjectStorage() {
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
    return (
        storageOpen === "loading"
            ? (
                <Center h="100%">
                    <Loader size={"lg"} />
                </Center>
            )
            : storageOpen
                ? (
                    <Box
                        p={"md"}
                        sx={theme => ({
                        })}>
                        <Stack>
                            <ObjectStorageInfo />
                            <ObjectPVCTable pvcs={pvcs} afterSubmit={updatePvc} />
                        </Stack>
                    </Box>
                )
                : (
                    <ObjectNotOpen openIt={() => {
                        openStorage("object")
                            .then(r => r.json())
                            .then(r => r.code === 200
                                ? Promise.resolve(r) : Promise.reject(r))
                            .then(r => noti("Success", "Success", "Open Storage"))
                            .catch(r => noti("Error", "Error", r.msg))
                            .finally(() => updateStorageOpen());
                    }} />
                )
    )
}

ObjectStorage.getLayout = commonLayout;

let hostColumns: MRT_ColumnDef<ObjectPVC>[] = [
    {
        header: "存储卷声明",
        accessorKey: "name"
    },
    {
        header: " 命名空间",
        accessorKey: "namespace",
    },
    {
        header: "对象数上限",
        accessorKey: "max_object"
    },
    {
        header: "容量上限",
        accessorKey: "max_size"
    },
    {
        header: "创建时间",
        accessorKey: "time"
    }
]

interface ObjectPVCTableInterface {
    pvcs: ObjectPVC[],
    afterSubmit: () => Promise<void>
}

function ObjectPVCTable({ pvcs, afterSubmit }: ObjectPVCTableInterface) {
    let router = useRouter();
    let form = useForm({
        initialValues: {
            name: "",
            namespace: "",
            autoscale: false,
            maxObject: 1,
            maxGbSize: 1
        },
        validate: {
            name: val => val.length === 0 ? "Empty Name" : null,
            namespace: val => val.length === 0 ? "Empty Namespace" : null,
            maxObject: val => (typeof val === "number" && val >= 1) ? null : "Invalid Value",
            maxGbSize: val => (typeof val === "number" && val >= 1) ? null : "Invalid Value",
        }
    });
    let [opened, { open, close }] = useDisclosure(false);
    let [editOpened, { close: editClose, open: editOpen }] = useDisclosure(false);
    let [deleteOpened, { close: deleteClose, open: deleteOpen }] = useDisclosure(false);
    let [currentRow, setCurrentRow] = useState<null | { name: string, namespace: string }>(null);
    useEffect(() => {
        form.reset();
    }, [opened]);
    return (
        <>
            <MantineReactTable
                columns={hostColumns}
                data={pvcs}
                enableColumnOrdering
                enableGlobalFilter
                enableBottomToolbar
                renderTopToolbarCustomActions={() => {
                    return (
                        <Button
                            variant="gradient"
                            gradient={{ from: 'teal', to: 'lime', deg: 105 }}
                            onClick={open}>
                            创建PVC
                        </Button>)
                }}
                onEditingRowSave={() => { }}
                enableEditing
                renderRowActions={({ row, table }) => (
                    <Box sx={{ display: 'flex', gap: '16px' }}>
                        <Tooltip withArrow position="right" label="删除">
                            <ActionIcon
                                color="red"
                                onClick={() => {
                                    setCurrentRow({
                                        name: row.original.name,
                                        namespace: row.original.namespace
                                    });
                                    deleteOpen();
                                }}>
                                <IconTrash />
                            </ActionIcon>
                        </Tooltip>
                        <Tooltip withArrow position="right" label="详情">
                            <ActionIcon
                                color="indigo"
                                onClick={() => {
                                    router.push(`/objectDetail?name=${row.original.name}&namespace=${row.original.namespace}`);
                                }}>
                                <IconListDetails />
                            </ActionIcon>
                        </Tooltip>
                    </Box>
                )}
            />
            <Modal
                opened={opened}
                onClose={close}
                centered
                title={<Title order={4}>新建PVC</Title>}
                radius={"lg"}
                withCloseButton={false}>
                <form onSubmit={form.onSubmit(values => {
                    createPvc({
                        name: values.name,
                        "X-type": "object",
                        autoscale: values.autoscale ? "true" : "false",
                        namespace: values.namespace,
                        maxgbsize: `${values.maxGbSize}`,
                        maxobject: `${values.maxObject}`
                    }).then(e => e.json())
                        .then((e) => {
                            return e.code === 200
                                ? Promise.resolve(e)
                                : Promise.reject(e);
                        })
                        .then(() => afterSubmit().then(() => {
                            close();
                            notifications.show({
                                id: v4(),
                                withCloseButton: true,
                                autoClose: 5000,
                                title: "Success",
                                message: "Create new PVC",
                                color: "green",
                                icon: <IconCheck />,
                                loading: false,
                            });
                        })).catch(e => {
                            close();
                            notifications.show({
                                id: v4(),
                                withCloseButton: true,
                                autoClose: 5000,
                                title: "Error",
                                message: e.msg,
                                color: "red",
                                icon: <IconX />,
                                loading: false,
                            });
                        });
                })}>
                    <Stack>
                        <TextInput
                            withAsterisk
                            label="存储卷声明"
                            {...form.getInputProps("name")} />
                        <TextInput
                            withAsterisk
                            label="命名空间"
                            {...form.getInputProps("namespace")} />
                        <NumberInput
                            min={1}
                            label="对象数上限"
                            {...form.getInputProps("maxObject")} />
                        <NumberInput
                            min={1}
                            label="容量上限"
                            {...form.getInputProps("maxGbSize")} />
                        <Switch
                            label="自动扩容"
                            {...form.getInputProps("autoscale")} />
                        <Button
                            variant="gradient"
                            gradient={{ from: 'teal', to: 'blue', deg: 105 }}
                            type="submit">
                            确定
                        </Button>
                    </Stack>
                </form>
            </Modal>
            <EditModal
                opened={editOpened}
                close={editClose}
                name={currentRow ? currentRow.name : ""}
                namespace={currentRow ? currentRow.namespace : ""}
                afterSubmit={afterSubmit}
            />
            <DeleteModal
                afterSubmit={afterSubmit}
                opened={deleteOpened}
                close={deleteClose}
                name={currentRow ? currentRow.name : ""}
                namespace={currentRow ? currentRow.namespace : ""}
            />
        </>
    )
}

interface EditModalInterface {
    name: string,
    namespace: string,
    afterSubmit: () => Promise<void>,
    opened: boolean,
    close: VoidFunction
}

function EditModal({ name, namespace, afterSubmit, opened, close }: EditModalInterface) {
    let form = useForm({
        initialValues: {
            volume: 1,
        },
        validate: {
            volume: val => (typeof val === "number" && val >= 1) ? null : "Invalid Value"
        }
    });
    useEffect(() => {
        form.reset();
    }, [opened]);
    return (
        <Modal
            opened={opened}
            onClose={close}
            centered
            title={<Title order={4}>编辑</Title>}
            size={"xs"}
            radius={"lg"}
            withCloseButton={false}>
            <form onSubmit={form.onSubmit(values => {
                updatePvcVolume({
                    name: name,
                    "X-type": "object",
                    namespace: namespace,
                    volume: `${values.volume}`
                }).then(e => e.json())
                    .then((e) => {
                        return e.code === 200
                            ? Promise.resolve(e)
                            : Promise.reject(e);
                    })
                    .then(() => afterSubmit().then(() => {
                        close();
                        notifications.show({
                            id: v4(),
                            withCloseButton: true,
                            autoClose: 5000,
                            title: "Success",
                            message: "Update PVC colume",
                            color: "green",
                            icon: <IconCheck />,
                            loading: false,
                        });
                    })).catch(e => {
                        close();
                        notifications.show({
                            id: v4(),
                            withCloseButton: true,
                            autoClose: 5000,
                            title: "Error",
                            message: e.msg,
                            color: "red",
                            icon: <IconX />,
                            loading: false,
                        });
                    });
            })}>
                <Stack>
                    <NumberInput
                        size="md"
                        min={1}
                        label="Volume"
                        {...form.getInputProps("volume")} />
                    <Button
                        variant="gradient"
                        gradient={{ from: 'teal', to: 'blue', deg: 105 }}
                        type="submit">
                        确定
                    </Button>
                </Stack>
            </form>
        </Modal>
    )
}

interface DeleteModalInterface {
    name: string,
    namespace: string,
    afterSubmit: () => Promise<void>,
    opened: boolean,
    close: VoidFunction
}

function DeleteModal({ name, namespace, afterSubmit, opened, close }: DeleteModalInterface) {
    return (
        <Modal
            opened={opened}
            onClose={close}
            centered
            size={"auto"}
            radius={"md"}
            withCloseButton={false}>
            <Stack>
                <Group position="apart" w="100%">
                    <Title order={5}>确定要删除此项PVC吗？</Title>
                    <Group spacing={"xs"}>
                        <Button size="xs" color="gray" onClick={close}>取消</Button>
                        <Button size="xs" color="red" onClick={() => {
                            deletePvc({ name, namespace, "X-type": "object" })
                                .then(e => e.json())
                                .then(e => e.code === 200 ? Promise.resolve(e) : Promise.reject(e))
                                .then(e => {
                                    notifications.show({
                                        id: v4(),
                                        withCloseButton: true,
                                        autoClose: 5000,
                                        title: "Success",
                                        message: "PVC deleted",
                                        color: "green",
                                        icon: <IconCheck />,
                                        loading: false,
                                    });
                                })
                                .catch(e => {
                                    notifications.show({
                                        id: v4(),
                                        withCloseButton: true,
                                        autoClose: 5000,
                                        title: "Error",
                                        message: e.msg,
                                        color: 'red',
                                        icon: <IconX />,
                                        loading: false,
                                        variant: "outline"
                                    });
                                })
                                .finally(() => {
                                    close();
                                    afterSubmit();
                                })
                        }}>删除</Button>
                    </Group>
                </Group>
            </Stack>
        </Modal>
    )
}

function ObjectNotOpen({ openIt }: { openIt: any }) {
    let theme = useMantineTheme();
    return (
        <Box h={"100%"}>
            <Center h={"100%"}>
                <Paper
                    radius={"xl"}
                    shadow="xl"
                    w={"60%"}
                    p={"xl"}>
                    <Stack>
                        <Box>
                            <Image
                                src={objectStorage}
                                height={150}
                                style={{
                                    float: "left",
                                    marginRight: 10,
                                    marginLeft: 10,
                                    marginTop: 10,
                                    clipPath: "circle(40%)"
                                }}
                                alt="BlockStorage" />
                            <Highlight
                                mt={30}
                                mr={30}
                                size={"md"}
                                highlight={['对象存储', '离散型数据', "桶", "对象桶"]}
                                highlightStyles={(theme) => ({
                                    backgroundImage: theme.fn.linearGradient(45, theme.colors.blue[5], theme.colors.grape[5]),
                                    fontWeight: 700,
                                    WebkitBackgroundClip: 'text',
                                    WebkitTextFillColor: 'transparent',
                                })}
                            >
                                对象存储是一种存储离散型数据的方法，结合了文件存储与块存储的优势。对象存储将存储数据称为对象，包含键Key，元数据MetaData，对象数据Data，用户可通过键Key访问对应的对象。对象的组织形式称为“桶”，在不同对象桶可以装载不同的对象。在 CubeUniverse 中，通过控制面板创建桶后，可以以 CubeUniverse RESTful 或 websocket 协议访问对象存储服务
                            </Highlight>
                        </Box>
                        <Center>
                            <Button
                                onClick={() => openIt()}
                                variant="gradient"
                                gradient={{ from: theme.colors.blue[5], to: theme.colors.grape[5], deg: 45 }}>开启块存储</Button>
                        </Center>
                    </Stack>
                </Paper>
            </Center>
        </Box>
    )
}

function ObjectStorageInfo() {
    let theme = useMantineTheme();
    // return <p>joe</p>
    return (
        <Center>
            <Paper
                p={"xl"}
                sx={theme => ({
                    borderRadius: 1,
                    boxShadow: "0px 2px 5px 1px rgb(113,88,219)",
                    border: "1px solid #c0e7ff"
                })}
            >
                <Stack>
                    <Box mr={100} ml={50} mt={20}>
                        <Image
                            height={150}
                            src={objectStorage}
                            style={{
                                float: "left",
                                marginRight: 30,
                                marginBottom: 20,
                                clipPath: "circle(40%)"
                            }}
                            alt="FileStorage" />
                            <Title
                            align="left"
                            order={3}
                            mb={"md"}
                        >对象存储</Title>
                        <Highlight
                            size={"md"}
                            highlight={['对象存储', '离散型数据', "桶", "对象桶"]}
                            highlightStyles={(theme) => ({
                                backgroundImage: theme.fn.linearGradient(45, theme.colors.blue[5], theme.colors.grape[5]),
                                fontWeight: 700,
                                WebkitBackgroundClip: 'text',
                                WebkitTextFillColor: 'transparent',
                            })}
                        >
                            对象存储是一种存储离散型数据的方法，结合了文件存储与块存储的优势。对象存储将存储数据称为对象，包含键Key，元数据MetaData，对象数据Data，用户可通过键Key访问对应的对象。对象的组织形式称为“桶”，在不同对象桶可以装载不同的对象。在 CubeUniverse 中，通过控制面板创建桶后，可以以 CubeUniverse RESTful 或 websocket 协议访问对象存储服务
                        </Highlight>
                    </Box>
                </Stack>
            </Paper>
        </Center>
    )
}