import { getUrls } from "@/api/url";
import { ShorturlData } from "@/types/ShorturlData";
import { FC, useEffect, useState } from "react";
//import ShorturlComponent from "./Shorturl";
import { AutoSizer, Column, Table, TableCellRenderer } from "react-virtualized";

type DeletedRowsState = { [key: string]: boolean };
type ShowDeleteConfirmationRowsState = { [key: string]: boolean };

type ShorturlsProps = {
  //data: ShorturlData[];
  shortUrls: ShorturlData[];
  setShortUrls: (shortUrls: ShorturlData[]) => void;
};

const ShorturlsComponent: FC<ShorturlsProps> = ({
  shortUrls,
  setShortUrls,
}) => {
  useEffect(() => {
    const fetchUrls = async () => {
      const data = await getUrls();
      setShortUrls(data);
    };
    fetchUrls();
  }, [setShortUrls]);

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
    navigator.clipboard.writeText(data.shortUrl);
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

  const handleConfirmDelete = (
    shouldDelete: boolean,
    rowData: ShorturlData,
  ) => {
    if (!rowData?.id) {
      throw new Error("Row data is missing");
    }

    if (shouldDelete) {
      markRowAsDeleted(rowData.id);
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
          <p className="no-interaction">⚠Delete? &nbsp;</p>
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
    return cellData
      ? `${cellData.toLocaleDateString()} ${cellData.toLocaleTimeString()}`
      : "-";
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
              cellDataGetter={({ rowData }) => rowData.useCount ?? 0}
            />
            <Column
              label="Last Used"
              dataKey="lastUsed"
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
