import {
  deleteUrl,
  getAuthenticatedWebsocket,
  getUrls,
  getUserMetadata,
} from "@/api/url";
import { ShorturlData, UserMetadata } from "@/types/ShorturlData";
import { FC, SetStateAction, useEffect, useState } from "react";
import { AutoSizer, Column, Table, TableCellRenderer } from "react-virtualized";
import { isDate } from "util/types";

type DeletedRowsState = { [key: string]: boolean };
type ShowDeleteConfirmationRowsState = { [key: string]: boolean };

type ShorturlsProps = {
  //data: ShorturlData[];
  shortUrls: ShorturlData[];
  setShortUrls: (value: SetStateAction<ShorturlData[]>) => void;
};

const ShorturlsComponent: FC<ShorturlsProps> = ({
  shortUrls,
  setShortUrls,
}) => {
  const [userMetadata, setUserMetadata] = useState<UserMetadata>();
  useEffect(() => {
    const fetchMetadata = async () => {
      const metadata = await getUserMetadata();
      setUserMetadata(metadata);
    };

    const fetchUrls = async () => {
      const data = await getUrls();
      setShortUrls(data);
    };
    fetchMetadata();
    fetchUrls();
  }, [setShortUrls]);

  let ws: WebSocket = null as any;
  useEffect(() => {
    // live updates
    const setupWebsocket = () => {
      ws = getAuthenticatedWebsocket();
      let authenticated = false;
      ws.onmessage = (event) => {
        try {
          console.log("Websocket message:", event.data);
          const data = JSON.parse(event.data);
          if (data.action === "auth") {
            authenticated = true;
          } else if (data.action === "created") {
            const shorturlData = data.data as ShorturlData;
            setShortUrls((prevState) => {
              // remove the same shorturl if it exists
              const newState = prevState.filter(
                (shorturl) => shorturl.id !== shorturlData.id,
              );
              return [shorturlData, ...newState];
            });
          } else if (data.action === "deleted") {
            const shorturlData = data.data as ShorturlData;
            setShortUrls((prevState) =>
              prevState.filter((shorturl) => shorturl.id !== shorturlData.id),
            );
          } else if (data.action === "updated") {
            const shorturlData = data.data as ShorturlData;
            setShortUrls((prevState) =>
              prevState.map((shorturl) =>
                shorturl.id === shorturlData.id ? shorturlData : shorturl,
              ),
            );
          }
        } catch (e) {
          console.error("Failed to parse websocket message", e);
        }
      };

      ws.onclose = (event) => {
        authenticated = false;
        console.log(
          "WebSocket closed with code",
          event.code,
          "and reason",
          event,
          "\nReconnecting in 5 seconds...",
        );
        setTimeout(setupWebsocket, 5 * 1000);
      };

      // heartbeat ping
      setInterval(() => {
        if (authenticated) {
          ws.send(
            JSON.stringify({
              action: "ping",
            }),
          );
        }
      }, 60 * 1000);
    };

    setupWebsocket();
    return () => {
      ws.close();
    };
    // websocket
  }, [setShortUrls, ws]);

  const [deletedRows, setDeletedRows] = useState<DeletedRowsState>({});
  const [showDeleteConfirmationRows, setShowDeleteConfirmationRows] =
    useState<ShowDeleteConfirmationRowsState>({});
  const [mostRecentShownDeleteRowId, setMostRecentShownHandleDeleteRowId] =
    useState("");

  const markRowAsDeleted = (rowId: string) => {
    setDeletedRows((prevState) => ({
      ...prevState,
      [rowId]: true,
    }));
  };

  const toggleShowDeleteConfirmation = (rowId: string) => {
    setShowDeleteConfirmationRows((prevState) => {
      if (mostRecentShownDeleteRowId && prevState[rowId]) {
        setMostRecentShownHandleDeleteRowId("");
      } else {
        setMostRecentShownHandleDeleteRowId(rowId);
      }

      return {
        ...prevState,
        [rowId]: !prevState[rowId],
      };
    });
  };

  const handleCopyShortUrl = (data: ShorturlData) => {
    navigator.clipboard.writeText(data.short_url);
  };

  const handleCopyUrl = (data: ShorturlData) => {
    navigator.clipboard.writeText(data.url);
  };

  const handleDelete = (data: ShorturlData) => {
    if (deletedRows[data.id]) return;
    console.log(`Delete ${data.id}`);
    if (mostRecentShownDeleteRowId) {
      toggleShowDeleteConfirmation(mostRecentShownDeleteRowId);
    }
    toggleShowDeleteConfirmation(data.id);
  };

  const handleConfirmDelete = async (
    shouldDelete: boolean,
    rowData: ShorturlData,
  ) => {
    if (!rowData?.id) {
      throw new Error("Row data is missing");
    }

    if (shouldDelete) {
      const didDelete = await deleteUrl(rowData.id);
      if (didDelete) {
        markRowAsDeleted(rowData.id);
        setShortUrls((prevState) =>
          prevState.filter((url) => url.id !== rowData.id),
        );

        const metadata = await getUserMetadata();
        setUserMetadata(metadata);
      }
    } else {
      toggleShowDeleteConfirmation(rowData?.id);
    }
  };

  const actionsCellRenderer: TableCellRenderer = ({ rowData }) => {
    return (
      <div hidden={deletedRows[rowData.id]}>
        <button
          type="button"
          className="action-delete"
          onClick={() => handleDelete(rowData as ShorturlData)}
          hidden={showDeleteConfirmationRows[rowData.id]}
        >
          ❌
        </button>
        <div hidden={!showDeleteConfirmationRows[rowData.id]}>
          <p className="no-interaction mobile-hide">⚠Delete? &nbsp;</p>
          <button
            type="button"
            className="destroy"
            onClick={() => handleConfirmDelete(true, rowData as ShorturlData)}
          >
            Yes
          </button>
          <p className="no-interaction">&nbsp;|&nbsp;</p>
          <button
            type="button"
            className="save"
            onClick={() => handleConfirmDelete(false, rowData as ShorturlData)}
          >
            No
          </button>
        </div>
      </div>
    );
  };

  const idCellRenderer: TableCellRenderer = ({ cellData, rowData }) => {
    return (
      <button
        type="button"
        className="copy"
        onClick={() => handleCopyShortUrl(rowData as ShorturlData)}
      >
        {cellData}
      </button>
    );
  };

  const lastUsedCellRenderer: TableCellRenderer = ({ cellData }) => {
    if (cellData) {
      const date = new Date(cellData);
      // remove seconds from the time string
      const time = date.toLocaleTimeString().split(":").slice(0, 2).join(":");
      return `${date.toLocaleDateString()} ${time}`;
    }
    return "-";
  };

  const urlCellRenderer: TableCellRenderer = ({ cellData, rowData }) => {
    return (
      <button
        type="button"
        className="copy"
        onClick={() => handleCopyUrl(rowData as ShorturlData)}
      >
        {cellData}
      </button>
    );
  };

  return (
    <div className="shorturls-list">
      <AutoSizer>
        {({ height, width }) => (
          <Table
            ref="Table"
            headerHeight={40}
            height={height}
            rowCount={shortUrls.length}
            rowHeight={48}
            rowGetter={({ index }) => shortUrls[index]}
            rowClassName={({ index }) =>
              index % 2 === 0 ? "even-row" : "odd-row"
            }
            width={width}
          >
            <Column
              label="ID"
              dataKey="id"
              width={width * 0.3}
              cellRenderer={idCellRenderer}
            />
            <Column
              label="URL"
              dataKey="url"
              width={width * 0.3}
              cellRenderer={urlCellRenderer}
            />
            <Column
              label="Uses"
              dataKey="uses"
              width={width * 0.05}
              cellDataGetter={({ rowData }) => rowData.use_count ?? 0}
            />
            <Column
              label="Last Used"
              dataKey="last_used"
              width={width * 0.15}
              cellRenderer={lastUsedCellRenderer}
            />
            <Column
              label="Actions"
              dataKey="actions"
              width={width * 0.2}
              cellRenderer={actionsCellRenderer}
            />
          </Table>
        )}
      </AutoSizer>
    </div>
  );
};

export default ShorturlsComponent;
