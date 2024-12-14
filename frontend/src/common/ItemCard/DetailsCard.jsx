import { useState } from 'react';
import { enqueueSnackbar } from 'notistack';

import dayjs from 'dayjs';
import { produce } from 'immer';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import { Card, CardMedia } from '@mui/material';
import relativeTime from 'dayjs/plugin/relativeTime';

import SimpleModal from '../SimpleModal';
import SharableGroups from '../SharableGroups';
import ImagePicker from '../ImagePicker/ImagePicker';
import DetailsCardItemContent from './ItemContent/DetailsCardItemContent';
import DetailsCardItemActions from './ItemContent/DetailsCardItemActions';
import { categoryActions } from '../../features/Categories/categoriesSlice';
import { maintenancePlanActions } from '../../features/MaintenancePlan/maintenanceSlice';
import { categoryItemDetailsActions } from '../../features/CategoryItemDetails/categoryItemDetailsSlice';
import { maintenancePlanItemActions } from '../../features/MaintenancePlanItemDetails/maintenancePlanItemSlice';

dayjs.extend(relativeTime);

export default function DetailsCard({ selectedItem, selectedImage, categoryMode = false }) {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const userID = localStorage.getItem('userID');

  const [editImgMode, setEditImgMode] = useState(null);
  const [openModal, setOpenModal] = useState(false);

  const handleOpenModal = () => setOpenModal(true);
  const handleCloseModal = () => setOpenModal(false);

  const handleUpload = (id, selectedImage) => {
    if (categoryMode) {
      dispatch(categoryItemDetailsActions.uploadImage({ id: id, selectedImage: selectedImage }));
    } else {
      dispatch(maintenancePlanItemActions.uploadImage({ id: id, selectedImage: selectedImage }));
    }
    enqueueSnackbar('New image upload successful.', {
      variant: 'success',
    });
    setEditImgMode(false);
  };

  const updateCollaborators = (sharableGroups) => {
    const newMembers = sharableGroups.map((v) => v.value);
    const draftSelectionDetails = produce(selectedItem, (draft) => {
      draft.updated_by = userID;
      draft.sharable_groups = newMembers;
      if (categoryMode) {
        draft.status = draft.status_name;
      } else {
        draft.maintenance_status = draft.maintenance_status_name;
      }
    });
    if (categoryMode) {
      dispatch(categoryActions.updateCategory(draftSelectionDetails));
      enqueueSnackbar('Updated collaborators for selected category.', {
        variant: 'success',
      });
      if (!newMembers.includes(userID)) {
        navigate('/');
      }
    } else {
      dispatch(maintenancePlanActions.updatePlan(draftSelectionDetails));
      enqueueSnackbar('Updated collaborators for selected maintenance plan.', {
        variant: 'success',
      });
      if (!newMembers.includes(userID)) {
        navigate('/');
      }
    }
    handleCloseModal();
  };

  return (
    <>
      <Card>
        <CardMedia sx={{ height: '10rem' }} image={selectedImage || '/blank_canvas.png'} />
        <DetailsCardItemContent selectedItem={selectedItem} categoryMode={categoryMode} />
        <DetailsCardItemActions
          selectedItem={selectedItem}
          handleOpenModal={handleOpenModal}
          setEditImgMode={setEditImgMode}
        />
      </Card>
      {openModal && (
        <SimpleModal
          title="Add sharable groups"
          subtitle="Assign collaborators."
          handleClose={handleCloseModal}
          maxSize="sm"
        >
          <SharableGroups
            handleSubmit={updateCollaborators}
            existingGroups={selectedItem?.sharable_groups || []}
            creator={selectedItem?.created_by}
          />
        </SimpleModal>
      )}
      {editImgMode && (
        <SimpleModal
          title="Assign image"
          subtitle="Assign image to the selected item."
          handleClose={() => setEditImgMode(false)}
          maxSize="xs"
        >
          <ImagePicker id={selectedItem.id} name={selectedItem.name} handleUpload={handleUpload} disableCancel />
        </SimpleModal>
      )}
    </>
  );
}
