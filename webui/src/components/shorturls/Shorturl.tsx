import { ShorturlData } from "@/types/ShorturlData";
import { FC, useState } from "react";

type ShorturlProps = {
  data?: ShorturlData;
  header?: boolean;
};

const ShorturlComponent: FC<ShorturlProps> = ({
  data = { id: "", shortUrl: "", url: "", uses: 0, lastUsed: new Date() },
  header = false,
}) => {
  const [deleted, setDeleted] = useState(false);
  const [showDeleteConfirmation, setShowDeleteConfirmation] = useState(false);

  let className = "shorturls-list-item";
  if (header) {
    className += " header";
  }

  const handleCopyShortUrl = () => {
    navigator.clipboard.writeText(data.shortUrl);
  };

  const handleCopyUrl = () => {
    navigator.clipboard.writeText(data.url);
  };

  const handleDelete = () => {
    if (deleted) return;
    console.log(`Delete ${data.id}`);
    setShowDeleteConfirmation(!showDeleteConfirmation);
  };

  const handleConfirmDelete = () => {
    setDeleted(true);
  };

  return header ? (
    <div className={className}>
      <div>ID</div>
      <div>URL</div>
      <div>Uses</div>
      <div>Last Used</div>
      <div>Actions</div>
    </div>
  ) : (
    <div key={data.id} className={className + (deleted ? " deleted" : "")}>
      <div>
        <button type="button" className="copy" onClick={handleCopyShortUrl}>
          {data.id}
        </button>
      </div>
      <div>
        <button type="button" className="copy" onClick={handleCopyUrl}>
          {data.url}
        </button>
      </div>
      <div>{data.uses}</div>
      <div>
        {`${data.lastUsed.toLocaleDateString()} ${data.lastUsed.toLocaleTimeString()}`}
      </div>
      <div hidden={deleted}>
        <button
          type="button"
          onClick={handleDelete}
          hidden={showDeleteConfirmation}
        >
          Delete
        </button>
        <div hidden={!showDeleteConfirmation}>
          Delete? &nbsp;
          <button
            type="button"
            className="destroy"
            onClick={handleConfirmDelete}
          >
            Yes
          </button>
          &nbsp;|&nbsp;
          <button type="button" className="save" onClick={handleConfirmDelete}>
            No
          </button>
        </div>
      </div>
    </div>
  );
};

export default ShorturlComponent;
